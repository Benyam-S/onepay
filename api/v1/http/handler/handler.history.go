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

// HandleGetUserHistory is a handler func that handles the request for retriving user's history per page
func (handler *UserAPIHandler) HandleGetUserHistory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

	userHistory, pageCount := handler.app.UserHistory(opUser.UserID, pagenation, viewBys...)

	output, _ := tools.MarshalIndent(HistoriesContainer{
		Result: userHistory, CurrentPage: pagenation, PageCount: pageCount}, "", "\t", format)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// HandleMarkHistoriesAsViewed is a handler func that handles the request for marking user's histories as viewed
func (handler *UserAPIHandler) HandleMarkHistoriesAsViewed(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	err := handler.app.MarkUserHistoriesAsViewed(opUser.UserID)

	if err != nil {
		output, _ := tools.MarshalIndent(ErrorBody{Error: err.Error()}, "", "\t", format)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
}
