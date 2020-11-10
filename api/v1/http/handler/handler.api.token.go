package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// HandleInitLoginApp is a handler func that handles a request for logging into the system using OnePay app
func (handler *UserAPIHandler) HandleInitLoginApp(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Get localization data from IP Geo location
	lb := new(entity.LocalizationBag)

	// Checking if the user exists
	opUser, err := handler.uService.FindUserAlsoWPhone(identifier, lb)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Frozen user checking
	if handler.dService.UserIsFrozen(opUser.UserID) {
		http.Error(w, entity.FrozenAccountError, http.StatusForbidden)
		return
	}

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
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
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

		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Check if the user has enabled two step verification
	userPreference, err := handler.uService.FindUserPreference(opUser.UserID)
	if err != nil || !userPreference.TwoStepVerification {
		var apiClient *api.Client
		apiClients, err := handler.uService.SearchAPIClient(opUser.UserID, entity.APIClientTypeInternal)
		if err != nil {
			newAPIClient := new(api.Client)
			newAPIClient.APPName = entity.APIClientAppNameInternal
			newAPIClient.Type = entity.APIClientTypeInternal
			err = handler.uService.AddAPIClient(newAPIClient, opUser)
			if err != nil {
				http.Error(w, entity.InternalAPIClientError, http.StatusInternalServerError)
				return
			}
			apiClient = newAPIClient
		} else {
			for _, client := range apiClients {
				if client.APPName == entity.APIClientAppNameInternal {
					apiClient = client
					break
				}
			}
		}

		// Frozen api client check
		if handler.dService.ClientIsFrozen(apiClient.APIKey) {
			http.Error(w, entity.FrozenAPIClientError, http.StatusForbidden)
			return
		}

		newAPIToken := new(api.Token)
		newAPIToken.Scopes = entity.ScopeAll
		err = handler.uService.AddAPIToken(newAPIToken, apiClient, opUser)
		if err != nil {
			http.Error(w, entity.APITokenError, http.StatusInternalServerError)
			return
		}

		output, _ := tools.MarshalIndent(map[string]string{"access_token": newAPIToken.AccessToken,
			"api_key": apiClient.APIKey, "type": "Bearer"}, "", "\t", format)
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	} else {

		// If the user has enabled two step verification the send otp
		otp := tools.GenerateOTP()
		smsNonce := uuid.Must(uuid.NewRandom())

		wd, _ := os.Getwd()
		dir := filepath.Join(wd, "./assets/messages", "/message.sms.otp.json")
		data, err1 := ioutil.ReadFile(dir)

		var messageSMS map[string][]string
		err2 := json.Unmarshal(data, &messageSMS)

		if err1 != nil || err2 != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// msg := messageSMS["message_body"][0] + otp + ". " + messageSMS["message_body"][1]
		// smsMessageID, err := tools.SendSMS(tools.OnlyPhoneNumber(opUser.PhoneNumber), msg)

		// if err != nil {
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }

		tempOutput, err1 := json.Marshal(opUser)
		err2 = tools.SetValue(handler.redisClient, smsNonce.String(), otp, time.Hour*6)
		err3 := tools.SetValue(handler.redisClient, otp+smsNonce.String(), string(tempOutput), time.Hour*6)
		if err1 != nil || err2 != nil || err3 != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Sending nonce to the client with the message ID, so it can be used to retrive the otp token
		output, _ := json.Marshal(map[string]string{"type": "OTP", "nonce": smsNonce.String(),
			"messageID": otp})
		w.WriteHeader(http.StatusOK)
		w.Write(output)

	}

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.PasswordFault+opUser.UserID)
}

// HandleVerifyLoginOTP is a handler func that handle a request for verifying otp token for login process
func (handler *UserAPIHandler) HandleVerifyLoginOTP(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")
	opUser := new(entity.User)

	// Analyzing nonce and otp
	if err := tools.AnalyzeKeyValuePair(handler.redisClient, nonce, otp); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	storedOPUser, err := tools.GetValue(handler.redisClient, otp+nonce)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid token used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Removing key value pair from the redis store
	tools.RemoveValues(handler.redisClient, otp+nonce)

	// unMarshaling user data
	json.Unmarshal([]byte(storedOPUser), opUser)

	var apiClient *api.Client
	apiClients, err := handler.uService.SearchAPIClient(opUser.UserID, entity.APIClientTypeInternal)
	if err != nil {
		newAPIClient := new(api.Client)
		newAPIClient.APPName = entity.APIClientAppNameInternal
		newAPIClient.Type = entity.APIClientTypeInternal
		err = handler.uService.AddAPIClient(newAPIClient, opUser)
		if err != nil {
			http.Error(w, entity.InternalAPIClientError, http.StatusInternalServerError)
			return
		}
		apiClient = newAPIClient
	} else {
		for _, client := range apiClients {
			if client.APPName == entity.APIClientAppNameInternal {
				apiClient = client
				break
			}
		}
	}

	// Frozen api client check
	if handler.dService.ClientIsFrozen(apiClient.APIKey) {
		http.Error(w, entity.FrozenAPIClientError, http.StatusForbidden)
		return
	}

	newAPIToken := new(api.Token)
	newAPIToken.Scopes = entity.ScopeAll
	err = handler.uService.AddAPIToken(newAPIToken, apiClient, opUser)
	if err != nil {
		http.Error(w, entity.APITokenError, http.StatusInternalServerError)
		return
	}

	output, _ := tools.MarshalIndent(map[string]string{"access_token": newAPIToken.AccessToken,
		"api_key": apiClient.APIKey, "type": "Bearer"}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleLogout is a handler func that handles a logout request
func (handler *UserAPIHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	apiToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Deactivating the api token
	apiToken.Deactivated = true
	handler.uService.UpdateAPIToken(apiToken)

}

// HandleRefreshAPITokenDE is a handler func that handles a request for refreshing the daily expiration date of
// a given api token
func (handler *UserAPIHandler) HandleRefreshAPITokenDE(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	apiToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Refershing the daily expiration date
	apiToken.DailyExpiration = time.Now().Unix()
	handler.uService.UpdateAPIToken(apiToken)

}

// HandleAuthentication is a handler func that handles a request for authenticating an api client
func (handler *UserAPIHandler) HandleAuthentication(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	apiKey := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	responseType := r.FormValue("response_type")
	scopesString := r.FormValue("scope")
	state := r.FormValue("state")

	emptyState, _ := regexp.MatchString(`^\s*$`, state)

	if emptyState {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "empty state value used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	if responseType != "code" {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "response type must be set to code"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	scopesSlice := strings.Split(scopesString, " ")
	for _, scope := range scopesSlice {
		if !api.ValidScope(scope) {
			output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid scope requested"}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}
	}

	// Frozen api client check
	if handler.dService.ClientIsFrozen(apiKey) {
		http.Error(w, entity.FrozenAPIClientError, http.StatusForbidden)
		return
	}

	apiClient, err := handler.uService.FindAPIClient(apiKey)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	if apiClient.CallBack != redirectURI {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "unregistred redirect uri used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	nonce := uuid.Must(uuid.NewRandom())
	scopes := strings.Join(scopesSlice, ", ")

	tempOutput, _ := json.Marshal(map[string]string{"api_key": apiKey, "scope": scopes, "state": state})
	err = tools.SetValue(handler.redisClient, nonce.String(), string(tempOutput), time.Hour*6)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// you should display a login page with nonce

	output, _ := tools.MarshalIndent(map[string]string{"nonce": nonce.String()}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleInitAuthorization is a handler func that initialize the request for authorizing a user
func (handler *UserAPIHandler) HandleInitAuthorization(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	nonce := r.FormValue("nonce")
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Checking for empty value
	if len(nonce) == 0 {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid nonce used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Retriving value from redis store
	storedDataS, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid nonce used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Get localization data from IP Geo location
	lb := new(entity.LocalizationBag)

	// Checking if the user exists
	opUser, err := handler.uService.FindUserAlsoWPhone(identifier, lb)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Frozen user checking
	if handler.dService.UserIsFrozen(opUser.UserID) {
		http.Error(w, entity.FrozenAccountError, http.StatusForbidden)
		return
	}

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
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
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

		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.InvalidPasswordOrIdentifierError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// you should display a check your scope page

	// updating the previously stored nonce data
	storedData := make(map[string]string)
	json.Unmarshal([]byte(storedDataS), &storedData)
	storedData["user_id"] = opUser.UserID

	tempOutput, _ := json.Marshal(storedData)
	err = tools.SetValue(handler.redisClient, nonce, string(tempOutput), time.Hour*6)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.PasswordFault+opUser.UserID)

	output, _ := tools.MarshalIndent(map[string]string{"nonce": nonce}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleFinishAuthorization is a handler func that finishes the authorization process
func (handler *UserAPIHandler) HandleFinishAuthorization(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	nonce := r.FormValue("nonce")
	authorized := r.FormValue("authorized")

	// Checking for empty value
	if len(nonce) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retriving value from redis store
	storedDataS, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if authorized != "true" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	code := uuid.Must(uuid.NewRandom()).String()

	err = tools.SetValue(handler.redisClient, code, nonce, time.Hour*6)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	storedData := make(map[string]string)
	json.Unmarshal([]byte(storedDataS), &storedData)

	output, _ := tools.MarshalIndent(map[string]string{"code": code, "state": storedData["state"]},
		"", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleCodeExchange is a handler func that enables an api client to exchange code for access token
func (handler *UserAPIHandler) HandleCodeExchange(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	code := r.FormValue("code")
	grantType := r.FormValue("grant_type")
	apiSecret := r.FormValue("client_secret")

	// Checking for empty value
	if len(code) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if grantType != "authorization_code" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retriving value from redis store
	storedNonce, err := tools.GetValue(handler.redisClient, code)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// getting the stored value using the nonce
	storedDataS, err := tools.GetValue(handler.redisClient, storedNonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	storedData := make(map[string]string)
	json.Unmarshal([]byte(storedDataS), &storedData)

	apiClient, err := handler.uService.FindAPIClient(storedData["api_key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if apiClient.APISecret != apiSecret {
		http.Error(w, "invalid client secret used", http.StatusBadRequest)
		return
	}

	opUser := new(entity.User)
	opUser.UserID = storedData["user_id"]

	newAPIToken := new(api.Token)
	newAPIToken.Scopes = storedData["scope"]

	err = handler.uService.AddAPIToken(newAPIToken, apiClient, opUser)
	if err != nil {
		http.Error(w, entity.APITokenError, http.StatusInternalServerError)
		return
	}

	// Removing temporary data
	tools.RemoveValues(handler.redisClient, code)
	tools.RemoveValues(handler.redisClient, storedNonce)

	output, _ := tools.MarshalIndent(map[string]string{"access_token": newAPIToken.AccessToken,
		"type": "Bearer"}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
