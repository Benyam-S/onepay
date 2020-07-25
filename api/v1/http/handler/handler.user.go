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
	"strings"
	"time"

	"github.com/gorilla/mux"

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
	app         *app.OnePay
	uService    user.IService
	dService    deleted.IService
	redisClient *redis.Client
}

// NewUserAPIHandler is a function that returns a new user api handler
func NewUserAPIHandler(commonApp *app.OnePay, userService user.IService, deletedService deleted.IService,
	redisClient *redis.Client) *UserAPIHandler {
	return &UserAPIHandler{app: commonApp, uService: userService, dService: deletedService,
		redisClient: redisClient}
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
	tools.RemoveValues(handler.redisClient, otp)

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

	// unmarshaling user data
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
		http.Error(w, "unable to add an internal api client", http.StatusInternalServerError)
		return
	}

	newAPIToken := new(api.Token)
	newAPIToken.Scopes = entity.ScopeAll
	err = handler.uService.AddAPIToken(newAPIToken, newAPIClient, newOPUser)
	if err != nil {
		http.Error(w, "unable to create an api token", http.StatusInternalServerError)
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

	// the format can a json or xml
	format := mux.Vars(r)["format"]

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

// HandleUpateProfile is a handler func that handles a request for updating user profile
func (handler *UserAPIHandler) HandleUpateProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	opUser.FirstName = r.FormValue("first_name")
	opUser.LastName = r.FormValue("last_name")
	opUser.Email = r.FormValue("email")
	opUser.PhoneNumber = r.FormValue("phone_number")

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

	tempFile, err := app.ClosingFile(opUser, userHistories, linkedAccounts)

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