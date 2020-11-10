package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// AccessTokenAuthentication is a middleware that validates whether a request contain a valid onepay access token
func (handler *UserAPIHandler) AccessTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, accessToken, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

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

		// No need for deactivating expired api tokens
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

// APITokenDEValidation is a middleware that checks whether an api token hasn't passed it daily expiration time
func (*UserAPIHandler) APITokenDEValidation(next http.HandlerFunc) http.HandlerFunc {

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

// Authorization is a middleware that authorize whether a given api token has a valid onepay user
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
			http.Error(w, entity.FrozenAccountError, http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("onepay_user"), opUser)
		r = r.Clone(ctx)

		// updating the api token for better user experience
		ipAddress, _ := tools.GetIP(r)
		apiToken.DeviceInfo = r.UserAgent()
		apiToken.IPAddress = ipAddress
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
			http.Error(w, "token scope is unauthorized for the request", http.StatusForbidden)
			return
		}

		next(w, r)
	}

}

// PasswordFaultHandler is a middleware that checks if the provided password in the request is valid or not.
// If it is invalid it will register it as a fault.
func (handler *UserAPIHandler) PasswordFaultHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		format := mux.Vars(r)["format"]
		password := r.FormValue("password")

		// checking for false attempts
		falseAttempts, _ := tools.GetValue(handler.redisClient, entity.PasswordFault+opUser.UserID)
		attempts, _ := strconv.ParseInt(falseAttempts, 0, 64)
		if attempts >= 5 {
			output, _ := tools.MarshalIndent(ErrorBody{Error: entity.TooManyAttemptsError}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
		opPassword, err := handler.uService.FindPassword(opUser.UserID)
		if err != nil {
			output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordError}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		// Comparing the hashed password with the given password
		hasedPassword, _ := base64.StdEncoding.DecodeString(opPassword.Password)
		err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password+opPassword.Salt))
		if err != nil {

			// registering fault
			tools.SetValue(handler.redisClient, entity.PasswordFault+opUser.UserID,
				fmt.Sprintf("%d", attempts+1), time.Hour*24)

			output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordError}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		// clearing user's false attempts
		tools.RemoveValues(handler.redisClient, entity.PasswordFault+opUser.UserID)
		next(w, r)
	}
}
