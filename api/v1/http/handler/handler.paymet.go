package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/app"

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
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.AmountParsingError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	moneyToken, err := handler.app.CreatePaymentToken(opUser.UserID, amount)

	if err != nil {

		// Whitelisting errors
		if err.Error() == entity.TransactionBaseLimitError ||
			err.Error() == entity.DailyTransactionLimitError {

			output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(output)
			return
		}

		// Any errors other than the above should be an internal server error
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusInternalServerError)
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

		// If error is any of the below then it will break out return bad request
		// else it will enter the default section so it can return internal server error
		switch err.Error() {
		// Whitelisting errors
		case entity.ReceiverNotFoundError:
		case entity.InvalidMoneyTokenError:
		case entity.ExpiredMoneyTokenError:
		case entity.TransactionBaseLimitError:
		case entity.DailyTransactionLimitError:
		case entity.TransactionWSelfError:
		case entity.InvalidMethodError:
		case entity.SenderNotFoundError:
		case entity.InsufficientBalanceError:
		default:
			// Any errors other than the above should be an internal server error
			output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(output)
			return
		}

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return

	}

}

// HandleGetPaymentInfo is a handler func that handles a request for getting payment info from the provided code
func (handler *UserAPIHandler) HandleGetPaymentInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	code := r.FormValue("code")
	receiverID := opUser.UserID

	moneyToken, err := handler.app.MoneyTokenService.FindMoneyToken(code)

	switch {
	case err != nil:
		err = errors.New(entity.InvalidMoneyTokenError)
	case !moneyToken.ExpirationDate.After(time.Now()):
		err = errors.New(entity.ExpiredMoneyTokenError)
	case !app.AboveTransactionBaseLimit(moneyToken.Amount):
		err = errors.New(entity.TransactionBaseLimitError)
	case moneyToken.SenderID == receiverID:
		err = errors.New(entity.TransactionWSelfError)
	case moneyToken.Method != entity.MethodPaymentQRCode:
		err = errors.New(entity.InvalidMethodError)
	}

	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(moneyToken, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}
