package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
)

// AccessTokenAuthentication is a middleware that validates a request contain a valid onepay access token
func (handler *UserAPIHandler) AccessTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		accessToken := r.FormValue("onepay_access_token")

		apiTokens, err := handler.uService.FindAPIToken(accessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		apiToken := apiTokens[0]

		if handler.uService.ValidateAPIToken(apiToken) != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Adding the api token to the context
		ctx := context.Background()
		ctx = context.WithValue(ctx, entity.Key("onepay_api_token"), apiToken)
		r = r.WithContext(ctx)

		next(w, r)

	}
}

// APITokenDEValidation is a middleware that checks an api token hasn't passed it daily expiration time
func (UserAPIHandler) APITokenDEValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		if apiToken.PastDailyExpiration() {
			output, _ := json.Marshal(map[string]string{"error": "access token has exceeded it daily expiration time"})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		next(w, r)
	}
}

// Authorization is a middleware that authorize a given api token has a valid onepay user
func (handler *UserAPIHandler) Authorization(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		opUser, err := handler.uService.FindUser(apiToken.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("onepay_user"), opUser)
		r = r.WithContext(ctx)

		// updating the api token for better user experience
		apiToken.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()
		handler.uService.UpdateAPIToken(apiToken)

		next(w, r)
	}
}
