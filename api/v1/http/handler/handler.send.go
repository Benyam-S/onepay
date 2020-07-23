package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// HandleSendMoneyViaQRCode is a handler func that handles a request for sending money via qr code
func (handler *UserAPIHandler) HandleSendMoneyViaQRCode(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]

	amountString := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	moneyToken, err := handler.app.SendViaQRCode(opUser.UserID, amount, handler.redisClient)

	if err != nil && err.Error() == entity.MoneyTokenCheckpointError {

		// requesting reload
		handler.app.Channel <- "reload_money_token"

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		return
	}

	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(CodeBody{Code: moneyToken.Code}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleSendMoneyViaOnePayID is a handler func that handles a request for sending money via onepay id
func (handler *UserAPIHandler) HandleSendMoneyViaOnePayID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	format := mux.Vars(r)["format"]

	receiverID := r.FormValue("receiver_id")
	amountString := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountString, 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.app.SendViaOnePayID(opUser.UserID, receiverID, amount, handler.redisClient)

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
