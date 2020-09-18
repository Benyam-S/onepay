package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/entity"
)

// WebsocketAccessTokenAuthentication is a middleware that validates whether a websocket request contain a valid onepay access token
func (handler *UserAPIHandler) WebsocketAccessTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := mux.Vars(r)["api_key"]
		accessToken := mux.Vars(r)["access_token"]

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
			http.Error(w, entity.FrozenAPIClientError, http.StatusForbidden)
			return
		}

		// Adding the api token to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, entity.Key("onepay_api_token"), apiToken)
		r = r.WithContext(ctx)

		next(w, r)

	}
}
