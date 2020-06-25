package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/session"
)

// SessionAuthentication is a middleware that validates a request cookie contain a valid onepay session value
func (UserHandler) SessionAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(os.Getenv("onepay_cookie_name"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		signedString := cookie.Value
		clientSession, err := session.Extract(signedString)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, entity.Key("onepay_client_session"), clientSession)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// SessionDEValidation is a middleware that checks a user client session hasn't passed it daily expiration time
func (UserHandler) SessionDEValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientSession, ok := ctx.Value(entity.Key("onepay_client_session")).(*session.ClientSession)

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return

		}

		if clientSession.PastDailyExpiration() {
			output, _ := json.Marshal(map[string]string{"error": "session has exceeded it daily expiration time"})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		next(w, r)
	}
}

// Authorization is a middleware that authorize a given session has a valid onepay user
func (handler *UserHandler) Authorization(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientSession, ok := ctx.Value(entity.Key("onepay_client_session")).(*session.ClientSession)

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return

		}

		serverSessions, err := handler.uservice.FindSession(clientSession.SessionID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// If you want you can notify the user
		if serverSessions[0].DeviceInfo != r.UserAgent() {
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}

		opUser, err := handler.uservice.FindUser(serverSessions[0].UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("onepay_user"), opUser)
		r = r.WithContext(ctx)

		handler.uservice.UpdateSession(serverSessions[0])

		next(w, r)
	}
}
