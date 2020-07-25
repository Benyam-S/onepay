package handler

import (
	"net/http"
	"strconv"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleCreatePaymentToken is a handler func that handles a request for creating a payment token
func (handler *UserAPIHandler) HandleCreatePaymentToken(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	amountString := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "amount parsing error"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	moneyToken, err := handler.app.CreatePaymentToken(opUser.UserID, amount)

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

// HandlePayViaQRCode is a handler func that handles a request for paying via qr code
func (handler *UserAPIHandler) HandlePayViaQRCode(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	code := r.FormValue("code")
	err := handler.app.PayViaQRCode(opUser.UserID, code, handler.redisClient)

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
