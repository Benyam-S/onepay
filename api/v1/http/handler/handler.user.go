package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

// UserAPIHandler is a type that defines a user handler for api client
type UserAPIHandler struct {
	uservice    user.IService
	redisClient *redis.Client
}

// NewUserAPIHandler is a function that returns a new user api handler
func NewUserAPIHandler(userService user.IService, redisClient *redis.Client) *UserAPIHandler {
	return &UserAPIHandler{uservice: userService, redisClient: redisClient}
}

// HandleFinishAddUser is a handler func that handles a request for adding password and constructing user account
// This is different from the client HandleFinishAddUser because it will return an api client at the end of the request
func (handler *UserAPIHandler) HandleFinishAddUser(w http.ResponseWriter, r *http.Request) {
	newOPPassword := new(entity.UserPassword)
	newOPUser := new(entity.User)

	nonce := r.FormValue("nonce")
	newOPPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vPassword")

	err := handler.uservice.VerifyUserPassword(newOPPassword, vPassword)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	storedOPUser, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": "invalid token used"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Removing key value pair from the redis store
	tools.RemoveValues(handler.redisClient, nonce)

	// unmarshaling user data
	json.Unmarshal([]byte(storedOPUser), newOPUser)

	err = handler.uservice.AddUser(newOPUser, newOPPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newAPIClient := new(api.Client)
	newAPIClient.APPName = entity.APIClientAppNameInternal
	newAPIClient.Type = entity.APIClientTypeInternal
	err = handler.uservice.AddAPIClient(newAPIClient, newOPUser)
	if err != nil {
		http.Error(w, "unable to add an internal api client", http.StatusInternalServerError)
		return
	}

	newAPIToken := new(api.Token)
	err = handler.uservice.AddAPIToken(newAPIToken, newAPIClient, newOPUser)
	if err != nil {
		http.Error(w, "unable to create an api token", http.StatusInternalServerError)
		return
	}

	signedString, _ := tools.GenerateToken([]byte(os.Getenv("onepay_secret_key")), newAPIToken)
	output, _ := json.Marshal(map[string]interface{}{"api_token": signedString})
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleInitLoginApp is a handler func that handles a request for logging into the system using OnePay app
func (handler *UserAPIHandler) HandleInitLoginApp(w http.ResponseWriter, r *http.Request) {
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Checking if the user exists
	opUser, err := handler.uservice.FindUser(identifier)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
	opPassword, err := handler.uservice.FindPassword(opUser.UserID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Comparing the hashed password with the given password
	hasedPassword, _ := base64.StdEncoding.DecodeString(opPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password+opPassword.Salt))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var apiClient *api.Client
	apiClients, err := handler.uservice.FindAPIClient(opUser.UserID, entity.APIClientTypeInternal)
	if err != nil {
		newAPIClient := new(api.Client)
		newAPIClient.APPName = entity.APIClientAppNameInternal
		newAPIClient.Type = entity.APIClientTypeInternal
		err = handler.uservice.AddAPIClient(newAPIClient, opUser)
		if err != nil {
			http.Error(w, "unable to add an internal api client", http.StatusInternalServerError)
			return
		}
		apiClient = newAPIClient
	} else {
		apiClient = apiClients[0]
	}

	newAPIToken := new(api.Token)
	err = handler.uservice.AddAPIToken(newAPIToken, apiClient, opUser)
	if err != nil {
		http.Error(w, "unable to create an api token", http.StatusInternalServerError)
		return
	}
	signedString, _ := tools.GenerateToken([]byte(os.Getenv("onepay_secret_key")), newAPIToken)
	output, _ := json.Marshal(map[string]interface{}{"api_token": signedString, "type": "Bearer"})
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
