package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetUserMoneyTokens is a handler func that returns all the money tokens created by the user
func (handler *UserAPIHandler) HandleGetUserMoneyTokens(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	moneyTokens := handler.app.MoneyTokenService.SearchMoneyToken(opUser.UserID)
	output, _ := tools.MarshalIndent(moneyTokens, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleRefreshMoneyTokens is a handler func that handles a request for refreshing a certain set of money tokens
func (handler *UserAPIHandler) HandleRefreshMoneyTokens(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	type MoneyTokenError struct {
		Code  string
		Error string
	}

	format := mux.Vars(r)["format"]
	codesString := r.FormValue("codes")

	refreshedMoneyTokens := make([]*entity.MoneyToken, 0)
	nonRefreshedMoneyTokens := make([]*MoneyTokenError, 0)

	codes := make([]string, 0)
	empty, _ := regexp.MatchString(`^\s*$`, codesString)
	if !empty {
		reg := regexp.MustCompile(`\s+`)
		codes = reg.Split(strings.TrimSpace(codesString), -1)
	}

	for _, code := range codes {
		err := handler.app.RefreshMoneyToken(code, opUser.UserID)
		if err != nil {
			errMT := new(MoneyTokenError)
			errMT.Code = code
			errMT.Error = err.Error()
			nonRefreshedMoneyTokens = append(nonRefreshedMoneyTokens, errMT)
		} else {
			moneyToken, err := handler.app.MoneyTokenService.FindMoneyToken(code)
			if err == nil {
				refreshedMoneyTokens = append(refreshedMoneyTokens, moneyToken)
			}
		}
	}

	// Meaning some of the codes hasn't been refreshed
	if len(nonRefreshedMoneyTokens) != 0 {

		type CRefreshedMoneyTokens struct {
			NonRefreshed []*MoneyTokenError
			Refreshed    []*entity.MoneyToken
		}

		output, _ := tools.MarshalIndent(CRefreshedMoneyTokens{
			NonRefreshed: nonRefreshedMoneyTokens,
			Refreshed:    refreshedMoneyTokens}, "", "\t", format)
		w.WriteHeader(http.StatusConflict)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(refreshedMoneyTokens, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return

}

// HandleReclaimMoneyTokens is a handler func that handles a request for reclaiming a certain set of money tokens
func (handler *UserAPIHandler) HandleReclaimMoneyTokens(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	codesString := r.FormValue("codes")

	errMap := make(entity.ErrMap)
	codes := make([]string, 0)
	empty, _ := regexp.MatchString(`^\s*$`, codesString)
	if !empty {
		reg := regexp.MustCompile(`\s+`)
		codes = reg.Split(strings.TrimSpace(codesString), -1)
	}

	for _, code := range codes {
		err := handler.app.ReclaimMoneyToken(code, opUser.UserID)
		if err != nil {
			if err.Error() == entity.WalletCheckpointError {
				// requesting reload
				handler.app.Channel <- "reload_wallet"
			}
			errMap[code] = err
		}
	}

	// Meaning some of the codes hasn't been reclaimed
	if len(errMap) != 0 {
		output, _ := tools.MarshalIndent(errMap.StringMap(), "", "\t", format)
		w.WriteHeader(http.StatusConflict)
		w.Write(output)
		return
	}

}

// HandleRemoveMoneyTokens is a handler func that handles a request for removing multiple money tokens of a certain user
func (handler *UserAPIHandler) HandleRemoveMoneyTokens(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	codesString := r.FormValue("codes")

	errMap := make(entity.ErrMap)
	codes := make([]string, 0)
	empty, _ := regexp.MatchString(`^\s*$`, codesString)
	if !empty {
		reg := regexp.MustCompile(`\s+`)
		codes = reg.Split(strings.TrimSpace(codesString), -1)
	}

	for _, code := range codes {
		moneyToken, err := handler.app.MoneyTokenService.FindMoneyToken(code)
		if err != nil {
			errMap[code] = err
			continue
		}

		// if the money token method is transaction qr code we have to reclaim it instead
		if moneyToken.Method == entity.MethodTransactionQRCode {
			err := handler.app.ReclaimMoneyToken(code, opUser.UserID)
			if err != nil {
				if err.Error() == entity.WalletCheckpointError {
					// requesting reload
					handler.app.Channel <- "reload_wallet"
				}
				errMap[code] = err
			}
		} else if moneyToken.Method == entity.MethodPaymentQRCode {
			err := handler.app.RemoveMoneyToken(code, opUser.UserID)
			if err != nil {
				errMap[code] = err
			}
		}

	}

	// Meaning some of the codes hasn't been deleted
	if len(errMap) != 0 {
		output, _ := tools.MarshalIndent(errMap.StringMap(), "", "\t", format)
		w.WriteHeader(http.StatusConflict)
		w.Write(output)
		return
	}
}
