package handler

import (
	"net/http"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
)

// HandleGetAccountProviders is a handler func that handles the request for getting all the
// account providers registered in the system
func (handler *UserAPIHandler) HandleGetAccountProviders(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	accountProviders := handler.apService.AllAccountProviders()

	output, _ := tools.MarshalIndent(accountProviders, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}

// HandleGetAccountProvider is a handler func that handles the request for getting detail information about the
// account provider
func (handler *UserAPIHandler) HandleGetAccountProvider(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]
	id := mux.Vars(r)["id"]

	accountProvider, err := handler.apService.FindAccountProvider(id)
	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: "account provider not found"}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	output, _ := tools.MarshalIndent(accountProvider, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	return
}
