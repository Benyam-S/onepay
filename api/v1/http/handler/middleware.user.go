package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// AccessTokenAuthentication is a middleware that validates a request contain a valid onepay access token
func (handler *UserAPIHandler) AccessTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, accessToken, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// accessToken := r.FormValue("onepay_access_token")

		apiToken, err := handler.uService.FindAPIToken(accessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// checking if the provided api client is similar with the api token's api key
		if apiToken.APIKey != apiKey {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if handler.uService.ValidateAPIToken(apiToken) != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Frozen api client checking
		if handler.dService.ClientIsFrozen(apiToken.APIKey) {
			http.Error(w, "api client has been frozen", http.StatusForbidden)
			return
		}

		// Adding the api token to the context
		ctx := r.Context()
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		format := mux.Vars(r)["format"]

		if apiToken.PastDailyExpiration() {
			output, _ := tools.MarshalIndent(ErrorBody{Error: "access token has exceeded it daily expiration time"},
				"", "\t", format)
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		opUser, err := handler.uService.FindUser(apiToken.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Frozen user checking
		if handler.dService.UserIsFrozen(opUser.UserID) {
			http.Error(w, "account has been frozen", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("onepay_user"), opUser)
		r = r.Clone(ctx)

		// updating the api token for better user experience
		apiToken.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()
		handler.uService.UpdateAPIToken(apiToken)

		next(w, r)
	}
}

// AuthenticateScope is a middleware that checks if the api token scope is compliant with the request
func (handler *UserAPIHandler) AuthenticateScope(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		apiToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		scopes := apiToken.GetScopes()
		requestScope, err := api.RequestScope(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		scopeFlage := false
		for _, scope := range scopes {
			if scope == requestScope {
				scopeFlage = true
				break
			}
		}

		if !scopeFlage {
			http.Error(w, "token scope is unathorized for the request", http.StatusForbidden)
			return
		}

		next(w, r)
	}

}
