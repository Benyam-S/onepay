package handler

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/tools"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
	"github.com/go-redis/redis"
)

// UserHandler is a type that defines a user handler for http client
type UserHandler struct {
	uService    user.IService
	redisClient *redis.Client
}

// NewUserHandler is a function that returns a new user handler
func NewUserHandler(userService user.IService, redisClient *redis.Client) *UserHandler {
	return &UserHandler{uService: userService, redisClient: redisClient}
}

// HandleInitAddUser is a handler func that handles a request for initiating adding new user
func (handler *UserHandler) HandleInitAddUser(w http.ResponseWriter, r *http.Request) {

	// In HandleAddUser you should not worry about receiving a profile picture since it a sign up page

	newOPUser := new(entity.User)
	newOPUser.FirstName = r.FormValue("first_name")
	newOPUser.LastName = r.FormValue("last_name")
	newOPUser.Email = r.FormValue("email")
	newOPUser.PhoneNumber = r.FormValue("phone_number")

	// validating user profile and cleaning up
	errMap := handler.uService.ValidateUserProfile(newOPUser)

	if errMap != nil {
		output, _ := json.Marshal(errMap.StringMap())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Generating OTP code
	otp := tools.GenerateOTP()
	// Generating unique key for identifing the OTP token
	smsNonce := uuid.Must(uuid.NewRandom())

	// Reading message body from our asset folder
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "./assets/messages", "/message.sms.json")
	data, err1 := ioutil.ReadFile(dir)

	var messageSMS map[string][]string
	err2 := json.Unmarshal(data, &messageSMS)

	if err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	msg := messageSMS["message_body"][0] + otp + ". " + messageSMS["message_body"][1]
	smsMessageID, err := tools.SendSMS(newOPUser.PhoneNumber, msg)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Saving all the data to a temporary database
	tempOutput, err1 := json.Marshal(newOPUser)
	err2 = tools.SetValue(handler.redisClient, smsNonce.String(), otp, time.Hour*6)
	err3 := tools.SetValue(handler.redisClient, otp+smsNonce.String(), string(tempOutput), time.Hour*6)
	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Sending nonce to the client with the message ID, so it can be used to retrive the otp token
	output, _ := json.Marshal(map[string]string{"nonce": smsNonce.String(), "messageID": smsMessageID})
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleVerifyOTP is a handler func that handle a request for verifying otp token
func (handler *UserHandler) HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")

	// Checking for empty value
	if len(otp) == 0 || len(nonce) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retriving value from redis store
	storedOTP, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Checking if the otp provided by request match the otp from the database
	if storedOTP != otp {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Removing key value pair from the redis store
	tools.RemoveValues(handler.redisClient, otp)

	// The nonce below is a key where the new user data is stored
	output, _ := json.Marshal(map[string]string{"nonce": otp + nonce})
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// HandleFinishAddUser is a handler func that handles a request for adding password and constructing user account
func (handler *UserHandler) HandleFinishAddUser(w http.ResponseWriter, r *http.Request) {

	newOPPassword := new(entity.UserPassword)
	newOPUser := new(entity.User)

	nonce := r.FormValue("nonce")
	newOPPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vPassword")

	err := handler.uService.VerifyUserPassword(newOPPassword, vPassword)
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

	err = handler.uService.AddUser(newOPUser, newOPPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := json.Marshal(newOPUser)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleLogin is a handler func that handles a request for logging into the system
func (handler *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	// Checking if the user exists
	opUser, err := handler.uService.FindUser(identifier)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
	opPassword, err := handler.uService.FindPassword(opUser.UserID)
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

	// Creating client side session and saving cookie
	opSession := session.Create(opUser.UserID)
	err = opSession.Save(w)
	if err != nil {
		http.Error(w, "unable to set cookie", http.StatusInternalServerError)
		return
	}

	// Creating server side session and saving to the system
	err = handler.uService.AddSession(opSession, opUser, r)
	if err != nil {
		http.Error(w, "unable to set session", http.StatusInternalServerError)
		return
	}

}

// HandleDashboard is a handler func that hanles a request for access the dashboard
func (handler *UserHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

}

// HandleLogout is a handler func that hanles a logout request
func (handler *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clientSession, ok := ctx.Value(entity.Key("onepay_client_session")).(*session.ClientSession)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	serverSession, err := handler.uService.FindSession(clientSession.SessionID)
	if err == nil {
		// Instead of deleting the session deactivate them for further use
		serverSession.Deactivated = true
		handler.uService.UpdateSession(serverSession)
	}

	// Removing client cookie
	clientSession.Remove(w)

}
