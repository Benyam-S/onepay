package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/entity"
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

		// Creating a new response write so it can be passed to the next handler
		newW := httptest.NewRecorder()

		// Adding the client session to the context
		ctx := context.Background()
		ctx = context.WithValue(ctx, entity.Key("onepay_client_session"), clientSession)
		r = r.WithContext(ctx)

		next(newW, r)

		// Obtaining all the header key value pairs from the previously created response writer
		for k, v := range newW.HeaderMap {
			w.Header()[k] = v
		}

		content := newW.Body.Bytes()

		// Refreshing session expiration time if we didn't want to logout
		// Meaning we didn't request Set-Cookie from the inside handlers
		if len(w.Header()["Set-Cookie"]) == 0 {
			clientSession.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()
			clientSession.UpdatedAt = time.Now().Unix()
			clientSession.Save(w)
		}

		w.WriteHeader(newW.Code)
		w.Write(content)

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

		serverSession, err := handler.uService.FindSession(clientSession.SessionID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Deactivated means the user has logged out from this session
		if serverSession.Deactivated {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// If you want you can notify the user
		if serverSession.DeviceInfo != r.UserAgent() {
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}

		opUser, err := handler.uService.FindUser(serverSession.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("onepay_user"), opUser)
		r = r.WithContext(ctx)

		handler.uService.UpdateSession(serverSession)

		next(w, r)
	}
}
