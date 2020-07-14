package api

import (
	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api/v1/http/handler"
	"github.com/Benyam-S/onepay/tools"
)

// Start is a function that start the provided api version
func Start(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/init", handler.HandleInitAddUser)
	router.HandleFunc("/api/v1/oauth/verify", handler.HandleVerifyOTP)
	router.HandleFunc("/api/v1/oauth/finish", handler.HandleFinishAddUser)

	router.HandleFunc("/api/v1/oauth/login", handler.HandleInitLoginApp)

	router.HandleFunc("/api/v1/oauth/logout", tools.MiddlewareFactory(handler.HandleLogout, handler.Authorization,
		handler.AccessTokenAuthentication))

	router.HandleFunc("/api/v1/oauth/user/password", tools.MiddlewareFactory(handler.HandleChangePassword, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile", tools.MiddlewareFactory(handler.HandleGetProfile, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile", tools.MiddlewareFactory(handler.HandleUpateProfile, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/{user_id}/pic", handler.HandleGetPhoto).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/pic", tools.MiddlewareFactory(handler.HandleUploadPhoto, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/pic", tools.MiddlewareFactory(handler.HandleRemovePhoto, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("DELETE")
}
