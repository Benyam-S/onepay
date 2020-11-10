package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/app"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleReceiveViaQRCode is a handler func that handles a request for receiving money via qr code
func (handler *UserAPIHandler) HandleReceiveViaQRCode(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	code := r.FormValue("code")

	// checking for false attempts
	falseAttempts, _ := tools.GetValue(handler.redisClient, entity.ReceiveFault+opUser.UserID)
	attempts, _ := strconv.ParseInt(falseAttempts, 0, 64)
	if attempts >= 5 {
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.TooManyAttemptsError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.app.ReceiveViaQRCode(opUser.UserID, code)

	if err != nil && err.Error() == entity.WalletCheckpointError {

		// requesting reload
		handler.app.Channel <- "reload_wallet"

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		return
	}

	if err != nil && err.Error() != entity.HistoryCheckpointError {

		// registering fault
		tools.SetValue(handler.redisClient, entity.ReceiveFault+opUser.UserID,
			fmt.Sprintf("%d", attempts+1), time.Hour*24)

		// If error is any of the below then it will break out return bad request
		// else it will enter the default section so it can return internal server error
		switch err.Error() {
		// Whitelisting errors
		case entity.ReceiverNotFoundError:
		case entity.InvalidMoneyTokenError:
		case entity.ExpiredMoneyTokenError:
		case entity.TransactionBaseLimitError:
		case entity.TransactionWSelfError:
		case entity.InvalidMethodError:
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

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.ReceiveFault+opUser.UserID)

}

// HandleGetReceiveInfo is a handler func that handles a request for getting receive amount info from the provided code
func (handler *UserAPIHandler) HandleGetReceiveInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	code := r.FormValue("code")

	// checking for false attempts
	falseAttempts, _ := tools.GetValue(handler.redisClient, entity.ReceiveFault+opUser.UserID)
	attempts, _ := strconv.ParseInt(falseAttempts, 0, 64)
	if attempts >= 5 {
		output, _ := tools.MarshalIndent(ErrorBody{Error: entity.TooManyAttemptsError}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	moneyToken, err := handler.app.MoneyTokenService.FindMoneyToken(code)

	switch {
	case err != nil:
		err = errors.New(entity.InvalidMoneyTokenError)
	case !moneyToken.ExpirationDate.After(time.Now()):
		err = errors.New(entity.ExpiredMoneyTokenError)
	case !app.AboveTransactionBaseLimit(moneyToken.Amount):
		err = errors.New(entity.TransactionBaseLimitError)
	case moneyToken.SenderID == opUser.UserID:
		err = errors.New(entity.TransactionWSelfError)
	case moneyToken.Method != entity.MethodTransactionQRCode:
		err = errors.New(entity.InvalidMethodError)
	}

	if err != nil {

		// registering fault
		tools.SetValue(handler.redisClient, entity.ReceiveFault+opUser.UserID,
			fmt.Sprintf("%d", attempts+1), time.Hour*24)

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(moneyToken, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.ReceiveFault+opUser.UserID)
	return
}
