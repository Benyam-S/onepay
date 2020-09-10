package api

import (
	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api/v1/http/handler"
	"github.com/Benyam-S/onepay/tools"
)

// Start is a function that start the provided api version
func Start(handler *handler.UserAPIHandler, router *mux.Router) {

	userRoutes(handler, router)
	apiTokenRoutes(handler, router)
	transactionRoutes(handler, router)
	walletNHistoryRoutes(handler, router)
	linkedAccountRoutes(handler, router)
	moneyTokenRoutes(handler, router)

}

// userRoutes is a function that defines all the routes for user profile handling
func userRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/register/init", handler.HandleInitAddUser)

	router.HandleFunc("/api/v1/oauth/user/register/verify", handler.HandleVerifyOTP)

	router.HandleFunc("/api/v1/oauth/user/register/finish.{format:json|xml}", handler.HandleFinishAddUser)

	router.HandleFunc("/api/v1/oauth/user", tools.MiddlewareFactory(handler.HandleDeleteUser, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("DELETE")

	router.HandleFunc("/api/v1/oauth/user/{user_id}/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUser, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetProfile, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleUpateProfile, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/{user_id}/profile/pic", handler.HandleGetPhoto).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile/pic.{format:json|xml}", tools.MiddlewareFactory(handler.HandleUploadPhoto, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/profile/pic", tools.MiddlewareFactory(handler.HandleRemovePhoto, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("DELETE")

	router.HandleFunc("/api/v1/oauth/user/password.{format:json|xml}", tools.MiddlewareFactory(handler.HandleChangePassword, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("PUT")

	/* ++++++++++++++++++++++++++++++++++++++++++ SESSION MANAGEMENT ++++++++++++++++++++++++++++++++++++++++++ */

	router.HandleFunc("/api/v1/oauth/user/session.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetActiveSessions, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/session", tools.MiddlewareFactory(handler.HandleDeactivateSession, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	/* ++++++++++++++++++++++++++++++++++++++++++++ FORGOT PASSWORD +++++++++++++++++++++++++++++++++++++++++++ */

	router.HandleFunc("/api/v1/user/password/rest/init.{format:json|xml}", tools.MiddlewareFactory(handler.HandleInitForgotPassword)).
		Methods("POST")

	router.HandleFunc("/api/v1/user/password/rest/finish/{nonce}", tools.MiddlewareFactory(handler.HandleFinishForgotPassword)).
		Methods("POST")
}

// tokenRoutes is a function that defines all the routes for handling api tokens
func apiTokenRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/authenticate.{format:json|xml}", handler.HandleAuthentication)

	router.HandleFunc("/api/v1/oauth/authorize/init.{format:json|xml}", handler.HandleInitAuthorization)

	router.HandleFunc("/api/v1/oauth/authorize/finish.{format:json|xml}", handler.HandleFinishAuthorization)

	router.HandleFunc("/api/v1/oauth/code/exchange.{format:json|xml}", handler.HandleCodeExchange)

	router.HandleFunc("/api/v1/oauth/login/app.{format:json|xml}", handler.HandleInitLoginApp)

	router.HandleFunc("/api/v1/oauth/refresh.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRefreshAPITokenDE, handler.Authorization,
		handler.AccessTokenAuthentication))

	router.HandleFunc("/api/v1/oauth/logout", tools.MiddlewareFactory(handler.HandleLogout, handler.Authorization,
		handler.AccessTokenAuthentication))

}

// transactionRoutes is a function that defines all the routes for handling a transaction
func transactionRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/send/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleSendMoneyViaQRCode,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/send/id.{format:json|xml}", tools.MiddlewareFactory(handler.HandleSendMoneyViaOnePayID,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/receive/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleReceiveViaQRCode,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/pay/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandlePayViaQRCode,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/pay/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleCreatePaymentToken,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")
}

// walletNHistoryRoutes is a function that defines all the routes for accessing user wallet and it's history
func walletNHistoryRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/wallet.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserWallet,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/wallet/recharge.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRechargeWallet,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/wallet/withdraw.{format:json|xml}", tools.MiddlewareFactory(handler.HandleWithdrawFromWallet,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/wallet/drain.{format:json|xml}", tools.MiddlewareFactory(handler.HandleDrainWallet,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/history.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserHistory,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")
}

// linkedAccountRoutes is a function that defines all the routes for accessing linked accounts of a certain user
func linkedAccountRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/init.{format:json|xml}", tools.MiddlewareFactory(handler.HandleInitLinkAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/finish", tools.MiddlewareFactory(handler.HandleFinishLinkAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserLinkedAccounts,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRemoveLinkedAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("DELETE")

}

// moneyTokenRoutes is a function that defines all the routes for accessing the money tokens of a certain user
func moneyTokenRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/moneytoken.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserMoneyTokens,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/moneytoken/refresh.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRefreshMoneyTokens,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/moneytoken/reclaim.{format:json|xml}", tools.MiddlewareFactory(handler.HandleReclaimMoneyTokens,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/moneytoken/remove.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRemoveMoneyTokens,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")
}
