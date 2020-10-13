package handler

import (
	"net/http"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/middleman"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetUserLinkedAccounts is a handler func that handles the request for getting the user's all linked accounts
func (handler *UserAPIHandler) HandleGetUserLinkedAccounts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	linkedAccounts := handler.app.GetUserLinkedAccounts(opUser.UserID)

	output, _ := tools.MarshalIndent(linkedAccounts, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleGetAccountInfo is a handler func that handles the request for getting linked account's account info
func (handler *UserAPIHandler) HandleGetAccountInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	linkedAccountID := mux.Vars(r)["id"]

	linkedAccount, err := handler.app.LinkedAccountService.FindLinkedAccount(linkedAccountID)
	if err != nil || linkedAccount.UserID != opUser.UserID {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "linked account not found"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	accountInfo, err := middleman.GetAccountInfo(linkedAccount.AccountID, linkedAccount.AccessToken)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "unable to fetch linked account info"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	accountInfo.AccountID = linkedAccount.AccountID
	accountInfo.AccountProviderID = linkedAccount.AccountProviderID

	output, _ := tools.MarshalIndent(accountInfo, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleInitLinkAccount is a handler func that initialize a request for linking external account to a user
func (handler *UserAPIHandler) HandleInitLinkAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	accountID := r.FormValue("account_id")
	accountProviderID := r.FormValue("account_provider_id")

	nonce, err := handler.app.AddLinkedAccount(opUser.UserID, accountID, accountProviderID, handler.redisClient)
	if err != nil {

		// Whitelisting errors
		// add the middleman errors later
		if err.Error() == "account provider not found" ||
			err.Error() == "account has already been linked to other onepay account" {
			output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := tools.MarshalIndent(map[string]string{"nonce": nonce}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleFinishLinkAccount is a handler func that finish the process of linking external account to a user
func (handler *UserAPIHandler) HandleFinishLinkAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")

	linkedAccount, err := handler.app.VerifyLinkedAccount(otp, nonce, handler.redisClient)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, _ := tools.MarshalIndent(linkedAccount, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// HandleRemoveLinkedAccount is a handler func that handles the request for unlinking an external account from onepay user's
func (handler *UserAPIHandler) HandleRemoveLinkedAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	linkedAccountID := r.FormValue("linked_account")

	linkedAccount, err := handler.app.RemoveLinkedAccount(linkedAccountID, opUser.UserID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Adding the removed linked account to trash
	handler.dService.AddLinkedAccountToTrash(linkedAccount)

	// cleaning the access token value so it can't be displayed
	linkedAccount.AccessToken = ""

	output, _ := tools.MarshalIndent(linkedAccount, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
