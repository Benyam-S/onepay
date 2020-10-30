package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/accountprovider"
	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/deleted"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/app"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserAPIHandler is a type that defines a user handler for api client
type UserAPIHandler struct {
	sync.Mutex
	app                  *app.OnePay
	uService             user.IService
	dService             deleted.IService
	apService            accountprovider.IService
	redisClient          *redis.Client
	upgrader             websocket.Upgrader
	activeSocketChannels map[string][]chan interface{}
}

// NewUserAPIHandler is a function that returns a new user api handler
func NewUserAPIHandler(commonApp *app.OnePay, userService user.IService, deletedService deleted.IService,
	accountProviderService accountprovider.IService, redisClient *redis.Client, upgrader websocket.Upgrader) *UserAPIHandler {
	return &UserAPIHandler{app: commonApp, uService: userService, dService: deletedService,
		apService: accountProviderService, redisClient: redisClient, upgrader: upgrader}
}

/* +++++++++++++++++++++++++++++++++++++++++++++ ADDING NEW USER +++++++++++++++++++++++++++++++++++++++++++++ */

// HandleInitAddUser is a handler func that handles a request for initiating adding new user
func (handler *UserAPIHandler) HandleInitAddUser(w http.ResponseWriter, r *http.Request) {

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
	// Generating unique key for identifying the OTP token
	smsNonce := uuid.Must(uuid.NewRandom())

	// Reading message body from our asset folder
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
	// smsMessageID, err := tools.SendSMS(tools.OnlyPhoneNumber(newOPUser.PhoneNumber), msg)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	// Saving all the data to a temporary database
	tempOutput, err1 := json.Marshal(newOPUser)
	err2 = tools.SetValue(handler.redisClient, smsNonce.String(), otp, time.Hour*6)
	err3 := tools.SetValue(handler.redisClient, otp+smsNonce.String(), string(tempOutput), time.Hour*6)
	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Sending nonce to the client with the message ID, so it can be used to retrive the otp token
	output, _ := json.Marshal(map[string]string{"nonce": smsNonce.String(), "messageID": otp})
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleVerifyOTP is a handler func that handle a request for verifying otp token
func (handler *UserAPIHandler) HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
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
	tools.RemoveValues(handler.redisClient, nonce)

	// The nonce below is a key where the new user data is stored
	output, _ := json.Marshal(map[string]string{"nonce": otp + nonce})
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// HandleFinishAddUser is a handler func that handles a request for adding password and constructing user account
// This is different from the client HandleFinishAddUser because it will return an api client at the end of the request
func (handler *UserAPIHandler) HandleFinishAddUser(w http.ResponseWriter, r *http.Request) {

	newOPPassword := new(entity.UserPassword)
	newOPUser := new(entity.User)

	format := mux.Vars(r)["format"]
	nonce := r.FormValue("nonce")
	newOPPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vPassword")

	err := handler.uService.VerifyUserPassword(newOPPassword, vPassword)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	storedOPUser, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid token used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Removing key value pair from the redis store
	tools.RemoveValues(handler.redisClient, nonce)

	// unMarshaling user data
	json.Unmarshal([]byte(storedOPUser), newOPUser)

	err = handler.uService.AddUser(newOPUser, newOPPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Adding wallet to the user
	newOPWallet := new(entity.UserWallet)
	newOPWallet.UserID = newOPUser.UserID
	err = handler.app.WalletService.AddWallet(newOPWallet)
	if err != nil {
		// This is cleaning up if the wallet is not created
		handler.uService.DeleteUser(newOPUser.UserID)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newAPIClient := new(api.Client)
	newAPIClient.APPName = entity.APIClientAppNameInternal
	newAPIClient.Type = entity.APIClientTypeInternal
	err = handler.uService.AddAPIClient(newAPIClient, newOPUser)
	if err != nil {
		http.Error(w, entity.InternalAPIClientError, http.StatusInternalServerError)
		return
	}

	newAPIToken := new(api.Token)
	newAPIToken.Scopes = entity.ScopeAll
	err = handler.uService.AddAPIToken(newAPIToken, newAPIClient, newOPUser)
	if err != nil {
		http.Error(w, entity.APITokenError, http.StatusInternalServerError)
		return
	}

	output, _ := tools.MarshalIndent(map[string]string{"access_token": newAPIToken.AccessToken,
		"type": "Bearer", "api_key": newAPIClient.APIKey}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

/* +++++++++++++++++++++++++++++++++++++++++++ GETTING PROFILE DATA +++++++++++++++++++++++++++++++++++++++++++ */

// HandleGetProfile is a handler func that handles a request for getting or viewing user's profile
func (handler *UserAPIHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// the format can be a json or xml
	format := mux.Vars(r)["format"]

	output, _ := tools.MarshalIndent(opUser, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleGetUser is a handler func that handles a request for getting or viewing user's profile
// This method can be a little controversial because it allow users to view other's profile,
// Solved the problem by deducting unnecessary values
func (handler *UserAPIHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// the format can be a json or xml
	format := mux.Vars(r)["format"]
	userID := mux.Vars(r)["user_id"]
	opUser, err := handler.uService.FindUser(userID)

	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "user not found"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	if handler.dService.UserIsFrozen(opUser.UserID) {
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.FrozenAccountError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Deducting unnecessary entries
	opUser.CreatedAt = time.Now()
	opUser.UpdatedAt = time.Now()
	opUser.PhoneNumber = ""
	opUser.Email = ""
	opUser.ProfilePic = ""

	output, _ := tools.MarshalIndent(opUser, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleGetPhoto is a handler func that handles a request for getting or viewing the user profile pic
func (handler *UserAPIHandler) HandleGetPhoto(w http.ResponseWriter, r *http.Request) {

	// This handler doesn't use ctx rather it uses the direct user_id embedded inside the query
	userID := mux.Vars(r)["user_id"]

	opUser, err := handler.uService.FindUser(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if opUser.ProfilePic != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "./assets/profilepics", opUser.ProfilePic)
		http.ServeFile(w, r, filePath)

	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

/* +++++++++++++++++++++++++++++++++++++++++++++ UPDATING PROFILE +++++++++++++++++++++++++++++++++++++++++++++ */

// HandleUpdateBasicInfo is a handler func that handles a request for updating user's profile basic information
func (handler *UserAPIHandler) HandleUpdateBasicInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	opUser.FirstName = r.FormValue("first_name")
	opUser.LastName = r.FormValue("last_name")

	errMap := handler.uService.ValidateUserProfile(opUser)

	if errMap != nil {
		output, _ := tools.MarshalIndent(errMap.StringMap(), "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.uService.UpdateUser(opUser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleInitUpdateEmail is a handler func that handles a request for updating user's email address
func (handler *UserAPIHandler) HandleInitUpdateEmail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	opUser.Email = r.FormValue("email")

	errMap := handler.uService.ValidateUserProfile(opUser)

	if errMap != nil && errMap["email"] != nil {
		output, _ := tools.MarshalIndent(
			ErrorBody{Error: errMap["email"].Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	otp := tools.GenerateOTP()
	emailNonce := uuid.Must(uuid.NewRandom())

	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "./assets/messages", "/message.email.verification.json")
	data, err1 := ioutil.ReadFile(dir)

	var messageEmail map[string][]string
	err2 := json.Unmarshal(data, &messageEmail)

	if err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	verificationLink := fmt.Sprintf("http://%s:%s/api/v1/user/profile/email/verify?nonce=%s&otp=%s",
		os.Getenv("domain_name"), os.Getenv("server_port"), emailNonce, otp)
	msg := messageEmail["message_body"][0] + opUser.UserID +
		messageEmail["message_body"][1] + verificationLink +
		". " + messageEmail["message_body"][2]
	err := tools.SendEmail(opUser.Email, "OnePay Email Verification", msg)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tempOutput, err1 := json.Marshal(opUser)
	err2 = tools.SetValue(handler.redisClient, emailNonce.String(), otp, time.Hour*24)
	err3 := tools.SetValue(handler.redisClient, otp+emailNonce.String(), string(tempOutput), time.Hour*24)
	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// type NonceContainer struct {
	// 	Nonce     string
	// 	MessageID string
	// }

	// output, _ := tools.MarshalIndent(
	// 	NonceContainer{Nonce: emailNonce.String(), MessageID: otp}, "", "\t", format)
	// w.WriteHeader(http.StatusOK)
	// w.Write(output)
}

// HandleVerifyUpdateEmail is a handler func that handle a request for verifying updated email address
func (handler *UserAPIHandler) HandleVerifyUpdateEmail(w http.ResponseWriter, r *http.Request) {
	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")
	updatedUser := new(entity.User)

	// Checking for empty value
	if len(otp) == 0 || len(nonce) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	storedOTP, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if storedOTP != otp {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tools.RemoveValues(handler.redisClient, nonce)

	storedOPUser, err := tools.GetValue(handler.redisClient, otp+nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tools.RemoveValues(handler.redisClient, otp+nonce)

	json.Unmarshal([]byte(storedOPUser), updatedUser)

	err = handler.uService.UpdateUserSingleValue(updatedUser.UserID, "email", updatedUser.Email)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// HandleInitUpdatePhone is a handler func that handles a request for updating user's phone number
func (handler *UserAPIHandler) HandleInitUpdatePhone(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	opUser.PhoneNumber = r.FormValue("phone_number")

	errMap := handler.uService.ValidateUserProfile(opUser)

	if errMap != nil && errMap["phone_number"] != nil {
		output, _ := tools.MarshalIndent(
			ErrorBody{Error: errMap["phone_number"].Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	otp := tools.GenerateOTP()
	smsNonce := uuid.Must(uuid.NewRandom())

	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "./assets/messages", "/message.sms.verification.json")
	data, err1 := ioutil.ReadFile(dir)

	var messageSMS map[string][]string
	err2 := json.Unmarshal(data, &messageSMS)

	if err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// msg := messageSMS["message_body"][0] + opUser.UserID +
	// 	messageSMS["message_body"][1] + otp +
	// 	". " + messageSMS["message_body"][2]
	// smsMessageID, err := tools.SendSMS(tools.OnlyPhoneNumber(opUser.PhoneNumber, msg)

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

	type NonceContainer struct {
		Nonce     string `json:"nonce"`
		MessageID string `json:"message_id"`
	}

	output, _ := tools.MarshalIndent(
		NonceContainer{Nonce: smsNonce.String(), MessageID: otp}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleVerifyUpdatePhone is a handler func that handle a request for verifying updated phone number
func (handler *UserAPIHandler) HandleVerifyUpdatePhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")
	updatedUser := new(entity.User)

	if len(otp) == 0 || len(nonce) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	storedOTP, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if storedOTP != otp {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tools.RemoveValues(handler.redisClient, nonce)

	storedOPUser, err := tools.GetValue(handler.redisClient, otp+nonce)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tools.RemoveValues(handler.redisClient, otp+nonce)

	json.Unmarshal([]byte(storedOPUser), updatedUser)

	err = handler.uService.UpdateUserSingleValue(updatedUser.UserID, "phone_number", updatedUser.PhoneNumber)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// HandleChangePassword is a handler func that handles a request for changing user passwords
func (handler *UserAPIHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	vPassword := r.FormValue("new_vPassword")

	opPassword, err := handler.uService.FindPassword(opUser.UserID)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	hasedPassword, _ := base64.StdEncoding.DecodeString(opPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(oldPassword+opPassword.Salt))
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid old password used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	if newPassword == oldPassword {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "new password is identical with the old password"},
			"", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	newOPPassword := new(entity.UserPassword)
	newOPPassword.UserID = opUser.UserID
	newOPPassword.Password = newPassword

	err = handler.uService.VerifyUserPassword(newOPPassword, vPassword)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err = handler.uService.UpdatePassword(newOPPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleUploadPhoto is a handler func that handles a request for uploading profile pic
func (handler *UserAPIHandler) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	// checking for multipart form data, the image has to be sent in multipart form data
	fm, fh, err := r.FormFile("profile_pic")
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
	defer fm.Close()

	// Reading the stream
	tempFile, _ := ioutil.ReadAll(fm)
	tempFileType := http.DetectContentType(tempFile)
	newBufferReader := bytes.NewBuffer(tempFile)

	// checking if the sent file is image
	if !strings.HasPrefix(tempFileType, "image") {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid format sent"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// checking the file sent doesn't exceed the size limit
	if fh.Size > 5000000 {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "image exceeds the file size limit, 5MB"},
			"", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	path, _ := os.Getwd()
	suffix := ""
	endPoint := 0

	for i := len(fh.Filename) - 1; i >= 0; i-- {
		if fh.Filename[i] == '.' {
			endPoint = i
			break
		}
	}

	for ; endPoint < len(fh.Filename); endPoint++ {
		suffix += string(fh.Filename[endPoint])
	}

	prevFileName := opUser.ProfilePic
	newFileName := fmt.Sprintf("asset_%s%s%s", opUser.UserID, tools.GenerateRandomString(3), suffix)
	for newFileName == prevFileName {
		newFileName = fmt.Sprintf("asset_%s%s%s", opUser.UserID, tools.GenerateRandomString(3), suffix)
	}

	path = filepath.Join(path, "./assets/profilepics", newFileName)

	out, _ := os.Create(path)
	defer out.Close()

	_, err = io.Copy(out, newBufferReader)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = handler.uService.UpdateUserSingleValue(opUser.UserID, "profile_pic", newFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if prevFileName != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "./assets/profilepics", prevFileName)
		tools.RemoveFile(filePath)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++ REMOVING PROFILE DATA ++++++++++++++++++++++++++++++++++++++++++ */

// HandleDeleteUser is a method that handles the request for deleting a user account
func (handler *UserAPIHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	password := r.FormValue("password")

	// checking for false attempts
	falseAttempts, _ := tools.GetValue(handler.redisClient, entity.PasswordFault+opUser.UserID)
	attempts, _ := strconv.ParseInt(falseAttempts, 0, 64)
	if attempts >= 5 {
		http.Error(w, entity.TooManyAttemptsError, http.StatusBadRequest)
		return
	}

	// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
	opPassword, err := handler.uService.FindPassword(opUser.UserID)
	if err != nil {
		http.Error(w, entity.InvalidPasswordError, http.StatusBadRequest)
		return
	}

	// Comparing the hashed password with the given password
	hasedPassword, _ := base64.StdEncoding.DecodeString(opPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password+opPassword.Salt))
	if err != nil {
		// registering fault
		tools.SetValue(handler.redisClient, entity.PasswordFault+opUser.UserID,
			fmt.Sprintf("%d", attempts+1), time.Hour*24)

		http.Error(w, entity.InvalidPasswordError, http.StatusBadRequest)
		return
	}

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.PasswordFault+opUser.UserID)

	userHistories, linkedAccounts, err := handler.app.InitDeleteOnePayAccount(opUser.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiClients, err := handler.uService.SearchAPIClient(opUser.UserID, entity.APIClientTypeUnfiltered)
	if err != nil {
		apiClients = []*api.Client{}
	}

	opUser, err = handler.uService.DeleteUser(opUser.UserID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Deleting the user's wallet
	handler.app.WalletService.DeleteWallet(opUser.UserID)

	// Removing all the user's linked accounts
	handler.app.LinkedAccountService.DeleteLinkedAccounts(opUser.UserID)

	// Adding deleted user to trash
	handler.dService.AddUserToTrash(opUser)

	for _, linkedAccount := range linkedAccounts {
		// Adding linked account to trash
		handler.dService.AddLinkedAccountToTrash(linkedAccount)
	}

	// Unfreezing user if it has been frozen
	handler.dService.UnfreezeUser(opUser.UserID)

	// Unfreezing api clients if any
	for _, apiClient := range apiClients {
		handler.dService.UnfreezeClient(apiClient.APIKey)
	}

	// Getting all the deleted linked accounts
	linkedAccounts = handler.dService.SearchDeletedLinkedAccounts("user_id", opUser.UserID)
	linkedAccountContainers := make([]*entity.LinkedAccountContainer, 0)
	for _, linkedAccount := range linkedAccounts {

		accountProvider, err := handler.apService.FindAccountProvider(linkedAccount.AccountProviderID)
		accountProviderName := accountProvider.Name
		if err != nil {
			accountProviderName = "account provider has been removed"
		}

		linkedAccountContainer := new(entity.LinkedAccountContainer)
		linkedAccountContainer.ID = linkedAccount.ID
		linkedAccountContainer.UserID = linkedAccount.UserID
		linkedAccountContainer.AccountID = linkedAccount.AccountID
		linkedAccountContainer.AccountProviderID = linkedAccount.AccountProviderID
		linkedAccountContainer.AccountProviderName = accountProviderName

		linkedAccountContainers = append(linkedAccountContainers, linkedAccountContainer)
	}

	tempFile, err := app.ClosingFile(opUser, userHistories, linkedAccountContainers)

	if err == nil {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "./assets/temp", tempFile)
		http.ServeFile(w, r, filePath)
		tools.RemoveFile(filePath)
	}

}

// HandleRemovePhoto is a handler func that handles the request for removing user's profile pic
func (handler *UserAPIHandler) HandleRemovePhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := handler.uService.UpdateUserSingleValue(opUser.UserID, "profile_pic", "")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if opUser.ProfilePic != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "./assets/profilepics", opUser.ProfilePic)
		tools.RemoveFile(filePath)
	}

}

/* ++++++++++++++++++++++++++++++++++++++++++ SESSION MANAGEMENT ++++++++++++++++++++++++++++++++++++++++++ */

// HandleGetActiveSessions is a handler func that handles the request for viewing all the user's active sessions
func (handler *UserAPIHandler) HandleGetActiveSessions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	type SessionContainer struct {
		ID         string
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeviceInfo string
		IPAddress  string
		Type       string
	}

	sessionContainers := make([]*SessionContainer, 0)

	sessions, err := handler.uService.SearchSession(opUser.UserID)
	if err != nil {
		sessions = []*session.ServerSession{}
	}

	for _, session := range sessions {
		if !session.Deactivated {
			sessionContainer := new(SessionContainer)
			sessionContainer.CreatedAt = session.CreatedAt
			sessionContainer.UpdatedAt = session.UpdatedAt
			sessionContainer.DeviceInfo = session.DeviceInfo
			sessionContainer.IPAddress = session.IPAddress
			sessionContainer.Type = entity.ClientTypeWeb
			sessionContainer.ID = session.SessionID

			sessionContainers = append(sessionContainers, sessionContainer)
		}
	}

	apiClients, err := handler.uService.SearchAPIClient(opUser.UserID, entity.APIClientTypeInternal)
	for _, apiClient := range apiClients {

		apiTokens, err := handler.uService.SearchAPIToken(apiClient.APIKey)
		if err != nil {
			apiTokens = []*api.Token{}
		}
		for _, apiToken := range apiTokens {
			err = handler.uService.ValidateAPIToken(apiToken)
			if err == nil && !apiToken.Deactivated {
				sessionContainer := new(SessionContainer)
				sessionContainer.CreatedAt = apiToken.CreatedAt
				sessionContainer.UpdatedAt = apiToken.UpdatedAt
				sessionContainer.IPAddress = apiToken.IPAddress
				sessionContainer.DeviceInfo = apiToken.DeviceInfo
				sessionContainer.Type = entity.APIClientTypeInternal
				sessionContainer.ID = apiToken.AccessToken

				if apiClient.Type == entity.APIClientTypeExternal {
					sessionContainer.DeviceInfo = apiClient.APPName
					sessionContainer.Type = entity.APIClientTypeExternal
				}

				sessionContainers = append(sessionContainers, sessionContainer)
			}
		}
	}

	output, _ := tools.MarshalIndent(sessionContainers, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleDeactivateSessions is a handler func that handles the request for deactivating user's active session
func (handler *UserAPIHandler) HandleDeactivateSessions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	currentAPIToken, ok := ctx.Value(entity.Key("onepay_api_token")).(*api.Token)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	idsString := r.FormValue("ids")

	ids := make([]string, 0)
	nonDeactivatedSessions := make([]string, 0)
	deactivatedSessions := make([]string, 0)

	empty, _ := regexp.MatchString(`^\s*$`, idsString)
	if !empty {
		reg := regexp.MustCompile(`\s+`)
		ids = reg.Split(strings.TrimSpace(idsString), -1)
	}

	for _, id := range ids {
		session, err1 := handler.uService.FindSession(id)
		apiToken, err2 := handler.uService.FindAPIToken(id)
		if err1 != nil && err2 != nil {
			nonDeactivatedSessions = append(nonDeactivatedSessions, id)
			continue
		}

		if session != nil {
			session.Deactivated = true
			err := handler.uService.UpdateSession(session)
			if err != nil {
				nonDeactivatedSessions = append(nonDeactivatedSessions, id)
				continue
			}
			deactivatedSessions = append(deactivatedSessions, id)

		} else if apiToken != nil {
			// Checking for current session deactivation
			if apiToken.UserID != opUser.UserID ||
				apiToken.AccessToken == currentAPIToken.AccessToken {
				nonDeactivatedSessions = append(nonDeactivatedSessions, id)
				continue
			}

			apiToken.Deactivated = true
			err := handler.uService.UpdateAPIToken(apiToken)
			if err != nil {
				nonDeactivatedSessions = append(nonDeactivatedSessions, id)
				continue
			}
			deactivatedSessions = append(deactivatedSessions, id)
		}
	}

	// Meaning some of the sessions hasn't been deactivated
	if len(nonDeactivatedSessions) != 0 {

		type SessionsStatusContainer struct {
			NonDeactivatedSessions []string
			DeactivatedSessions    []string
		}

		output, _ := tools.MarshalIndent(
			SessionsStatusContainer{
				DeactivatedSessions:    deactivatedSessions,
				NonDeactivatedSessions: nonDeactivatedSessions,
			}, "", "\t", format)
		w.WriteHeader(http.StatusConflict)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(deactivatedSessions, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

/* ++++++++++++++++++++++++++++++++++++++++++++ FORGOT PASSWORD +++++++++++++++++++++++++++++++++++++++++++ */

// HandleInitForgotPassword is a handler func that initiate the forgot password process
func (handler *UserAPIHandler) HandleInitForgotPassword(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	method := r.FormValue("method")
	identifier := r.FormValue("identifier")

	// Get localization data from IP Geo location
	lb := new(entity.LocalizationBag)

	opUser, err := handler.uService.FindUserAlsoWPhone(identifier, lb)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "invalid identifier used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	if method == "email" {

		nonce := uuid.Must(uuid.NewRandom())

		// Reading message body from our asset folder
		wd, _ := os.Getwd()
		dir := filepath.Join(wd, "./assets/messages", "/message.email.rest.json")
		data, err1 := ioutil.ReadFile(dir)

		var messageEmail map[string][]string
		err2 := json.Unmarshal(data, &messageEmail)

		if err1 != nil || err2 != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		restPath := os.Getenv("domain_name") + ":" + os.Getenv("server_port") + "/user/password/rest/finish/" + nonce.String()
		msg := messageEmail["message_body"][0] + restPath + ". " + messageEmail["message_body"][1]
		err := tools.SendEmail(opUser.Email, "Rest OnePay User Password", msg)

		if err != nil {
			output, _ := tools.MarshalIndent(ErrorBody{Error: "unable to send message"}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		err = tools.SetValue(handler.redisClient, nonce.String(), opUser.UserID, time.Hour*6)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	} else if method == "phone_number" {

		nonce := uuid.Must(uuid.NewRandom())

		// Reading message body from our asset folder
		wd, _ := os.Getwd()
		dir := filepath.Join(wd, "./assets/messages", "/message.sms.rest.json")
		data, err1 := ioutil.ReadFile(dir)

		var messageSMS map[string][]string
		err2 := json.Unmarshal(data, &messageSMS)

		if err1 != nil || err2 != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		restPath := os.Getenv("domain_name") + ":" + os.Getenv("server_port") + "/user/password/rest/finish/" + nonce.String()
		msg := messageSMS["message_body"][0] + restPath + ". " + messageSMS["message_body"][1]
		_, err := tools.SendSMS(tools.OnlyPhoneNumber(opUser.PhoneNumber), msg)

		if err != nil {
			output, _ := tools.MarshalIndent(ErrorBody{Error: "unable to send message"}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		err = tools.SetValue(handler.redisClient, nonce.String(), opUser.UserID, time.Hour*6)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	} else {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "unknown method used"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}

// HandleFinishForgotPassword is a handler func that finish the forgot password process
func (handler *UserAPIHandler) HandleFinishForgotPassword(w http.ResponseWriter, r *http.Request) {

	newOPPassword := new(entity.UserPassword)

	nonce := mux.Vars(r)["nonce"]
	newOPPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vPassword")

	err := handler.uService.VerifyUserPassword(newOPPassword, vPassword)
	if err != nil {
		output, _ := json.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	storedOPUserID, err := tools.GetValue(handler.redisClient, nonce)
	if err != nil {
		output, _ := json.MarshalIndent(ErrorBody{Error: "invalid token used"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Removing key value pair from the redis store
	tools.RemoveValues(handler.redisClient, nonce)

	opUser, err := handler.uService.FindUser(storedOPUserID)
	if err != nil {
		output, _ := json.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	newOPPassword.UserID = opUser.UserID
	err = handler.uService.UpdatePassword(newOPPassword)
	if err != nil {
		output, _ := json.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		return
	}

}
