package api

import (
	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api/v1/http/handler"
)

// Start is a function that start the provided api version
func Start(handler *handler.UserAPIHandler, router *mux.Router) {
	router.HandleFunc("/api/v1/oauth/finish", handler.HandleFinishAddUser)
	router.HandleFunc("/api/v1/oauth/login", handler.HandleInitLoginApp)
}
