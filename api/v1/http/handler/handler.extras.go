package handler

import (
	"net/http"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetCurrencyRates is a handler func that gets the current currency rate for the provided base currency
func (handler *UserAPIHandler) HandleGetCurrencyRates(w http.ResponseWriter, r *http.Request) {

	format := mux.Vars(r)["format"]
	base := mux.Vars(r)["base"]

	rates, err := handler.app.GetCurrencyRates(base)
	if err != nil {
		if err.Error() == "unknown base currency" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := tools.MarshalIndent(rates, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleGetNotifications is a handler func that retrive any new changes occuring with in the recent 24 hours
func (handler *UserAPIHandler) HandleGetNotifications(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	newHistories := make([]*entity.UserHistory, 0)
	baseTime := time.Now().Add(time.Hour * time.Duration(-24))
	outDated := false

	for i := 0; true; i++ {
		histories, _ := handler.app.UserHistory(opUser.UserID, int64(i), "all")
		for _, history := range histories {
			// ReceivedAt is used since SentAt can be out-dated
			if history.ReceivedAt.Before(baseTime) {
				outDated = true
				break
			}

			if (history.SenderID == opUser.UserID && history.SenderSeen) ||
				(history.ReceiverID == opUser.UserID && history.ReceiverSeen) {
				continue
			}

			newHistories = append(newHistories, history)
		}

		if outDated {
			break
		}
	}

	output, _ := tools.MarshalIndent(NotificationContainer{Histories: newHistories}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
