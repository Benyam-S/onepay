package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetUserLinkedAccounts is a handler func that handles the request for getting the user's all linked accounts
func (handler *UserAPIHandler) HandleGetUserLinkedAccounts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]

	linkedAccounts := make([]*LinkedAccountBody, 0)
	linkedAccountsMap := handler.app.GetUserLinkedAccounts(opUser.UserID)

	for linkedAccount, accountInfo := range linkedAccountsMap {

		linkedAccountBody := new(LinkedAccountBody)
		linkedAccountBody.AccountID = linkedAccount.AccountID
		linkedAccountBody.AccountProvider = linkedAccount.AccountProvider
		linkedAccountBody.ID = linkedAccount.ID
		linkedAccountBody.Amount = accountInfo.Amount

		linkedAccounts = append(linkedAccounts, linkedAccountBody)
	}

	output, _ := tools.MarshalIndent(linkedAccounts, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleRechargeWallet is a handler func that handles the request for recharging user wallet
func (handler *UserAPIHandler) HandleRechargeWallet(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]
	linkedAccountID := r.FormValue("linked_account")
	amountString := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.app.RechargeWallet(opUser.UserID, linkedAccountID, amount)
	if err != nil && err.Error() == entity.WalletCheckpointError {

		// requesting reload
		handler.app.Channel <- "reload_wallet"

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		return
	}

	if err != nil && err.Error() != entity.HistoryCheckpointError {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}

// HandleWithdrawFromWallet is a handler func that handles the request for withdrawing money from user's wallet
func (handler *UserAPIHandler) HandleWithdrawFromWallet(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]
	linkedAccountID := r.FormValue("linked_account")
	amountString := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.app.WithdrawFromWallet(opUser.UserID, linkedAccountID, amount)
	if err != nil && err.Error() == entity.WalletCheckpointError {

		// requesting reload
		handler.app.Channel <- "reload_wallet"

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		return
	}

	if err != nil && err.Error() != entity.HistoryCheckpointError {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}

// HandleInitLinkAccount is a handler func that initialize a request for linking external account to a user
func (handler *UserAPIHandler) HandleInitLinkAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]
	accountID := r.FormValue("account_id")
	accountProvider := r.FormValue("account_provider")
	emptyAccountID, _ := regexp.MatchString(`^\s*$`, accountID)
	emptyAccountProvider, _ := regexp.MatchString(`^\s*$`, accountProvider)

	if emptyAccountID || emptyAccountProvider {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	nonce, err := handler.app.AddLinkedAccount(opUser.UserID, accountID, accountProvider, handler.redisClient)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	otp := r.FormValue("otp")
	nonce := r.FormValue("nonce")

	err := handler.app.VerifyLinkedAccount(otp, nonce, handler.redisClient)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

}

// HandleRemoveLinkedAccount is a handler func that handles the request for unlinking an external account from onepay user's
func (handler *UserAPIHandler) HandleRemoveLinkedAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
