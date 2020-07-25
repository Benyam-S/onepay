package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleReceiveViaQRCode is a handler func that handles a request for receivng money via qr code
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
		output, _ := tools.MarshalIndent(ErrorBody{Error: "too many attempts try after 24 hours"}, "", "\t", format)
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

		// registring fault
		tools.SetValue(handler.redisClient, entity.ReceiveFault+opUser.UserID,
			fmt.Sprintf("%d", attempts+1), time.Hour*24)

		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// clearing user's false attempts
	tools.RemoveValues(handler.redisClient, entity.ReceiveFault+opUser.UserID)

}
