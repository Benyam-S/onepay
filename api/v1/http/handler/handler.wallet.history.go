package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetUserWallet is a handler func that handles the request for getting the user's wallet
func (handler *UserAPIHandler) HandleGetUserWallet(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]

	opWallet, err := handler.app.WalletService.FindWallet(opUser.UserID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, _ := tools.MarshalIndent(opWallet, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleDrainWallet is a handler func that handles the request for draining/emptying a wallet
func (handler *UserAPIHandler) HandleDrainWallet(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]
	linkedAccountID := r.FormValue("linked_account")

	err := handler.app.DrainWallet(opUser.UserID, linkedAccountID)
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

// HandleGetUserHistory is a handler func that handles the request for retriving user's hitory per page
func (handler *UserAPIHandler) HandleGetUserHistory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]
	pageString := r.FormValue("page")
	viewBysString := r.FormValue("view_bys")
	viewBys := make([]string, 0)

	pagenation, _ := strconv.ParseInt(pageString, 0, 64)
	empty, _ := regexp.MatchString(`^\s*$`, viewBysString)
	if !empty {
		viewBys = strings.Split(strings.TrimSpace(viewBysString), " ")
	}

	// if no view by then return all
	if len(viewBys) == 0 {
		viewBys = append(viewBys, "all")
	}

	userHistory := handler.app.UserHistory(opUser.UserID, pagenation, viewBys...)

	output, _ := tools.MarshalIndent(userHistory, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
