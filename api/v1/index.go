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
	websocketRoutes(handler, router)

}

// userRoutes is a function that defines all the routes for user profile handling
func userRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/register/init", handler.HandleInitAddUser)

	router.HandleFunc("/api/v1/oauth/user/register/verify", handler.HandleVerifyAddUserOTP)

	router.HandleFunc("/api/v1/oauth/user/register/finish.{format:json|xml}", handler.HandleFinishAddUser)

	router.HandleFunc("/api/v1/oauth/user.{format:json|xml}", tools.MiddlewareFactory(handler.HandleDeleteUser,
		handler.PasswordFaultHandler, handler.Authorization, handler.AccessTokenAuthentication)).Methods("DELETE")

	router.HandleFunc("/api/v1/oauth/user/{user_id}/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUser, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetProfile, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	// router.HandleFunc("/api/v1/oauth/user/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleUpdateProfile, handler.Authorization,
	// 	handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile.{format:json|xml}", tools.MiddlewareFactory(handler.HandleUpdateBasicInfo, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile/preference", tools.MiddlewareFactory(handler.HandleUpdateUserPreference, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile/phonenumber.{format:json|xml}", tools.MiddlewareFactory(handler.HandleInitUpdatePhone, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile/phonenumber/verify", tools.MiddlewareFactory(handler.HandleVerifyUpdatePhone, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/user/profile/email.{format:json|xml}", tools.MiddlewareFactory(handler.HandleInitUpdateEmail, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	// Since we are using link for verifying the method should be get and also it doesn't use Oauth
	router.HandleFunc("/api/v1/user/profile/email/verify", handler.HandleVerifyUpdateEmail).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/{user_id}/profile/pic", handler.HandleGetPhoto).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/statement/{id}", handler.HandleGetAccountStatement).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/profile/pic.{format:json|xml}", tools.MiddlewareFactory(handler.HandleUploadPhoto, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/profile/pic", tools.MiddlewareFactory(handler.HandleRemovePhoto, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("DELETE")

	router.HandleFunc("/api/v1/oauth/user/password.{format:json|xml}", tools.MiddlewareFactory(handler.HandleChangePassword, handler.Authorization,
		handler.AccessTokenAuthentication)).Methods("PUT")

	/* ++++++++++++++++++++++++++++++++++++++++++ SESSION MANAGEMENT ++++++++++++++++++++++++++++++++++++++++++ */

	router.HandleFunc("/api/v1/oauth/user/session.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetActiveSessions, handler.Authorization,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/session.{format:json|xml}", tools.MiddlewareFactory(handler.HandleDeactivateSessions, handler.Authorization,
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

	router.HandleFunc("/api/v1/oauth/login/app/verify.{format:json|xml}", handler.HandleVerifyLoginOTP)

	router.HandleFunc("/api/v1/oauth/refresh.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRefreshAPITokenDE,
		handler.PasswordFaultHandler, handler.Authorization, handler.AccessTokenAuthentication))

	router.HandleFunc("/api/v1/oauth/logout", tools.MiddlewareFactory(handler.HandleLogout, handler.Authorization,
		handler.AccessTokenAuthentication))

	router.HandleFunc("/api/v1/oauth/resend", handler.HandleResendMessage).Methods("POST")
}

// transactionRoutes is a function that defines all the routes for handling a transaction
func transactionRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/send/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleSendMoneyViaQRCode,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/send/id.{format:json|xml}", tools.MiddlewareFactory(handler.HandleSendMoneyViaOnePayID,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/receive/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetReceiveInfo,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/receive/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleReceiveViaQRCode,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/pay/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetPaymentInfo,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/pay/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandlePayViaQRCode,
		handler.Authorization, handler.APITokenDEValidation,
		handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

	router.HandleFunc("/api/v1/oauth/pay/code.{format:json|xml}", tools.MiddlewareFactory(handler.HandleCreatePaymentToken,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")
}

// walletNHistoryRoutes is a function that defines all the routes for accessing user wallet and it's history
func walletNHistoryRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/wallet.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserWallet,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/wallet.{format:json|xml}", tools.MiddlewareFactory(handler.HandleMarkWalletAsViewed,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")

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

	router.HandleFunc("/api/v1/oauth/user/history.{format:json|xml}", tools.MiddlewareFactory(handler.HandleMarkHistoriesAsViewed,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("PUT")
}

// linkedAccountRoutes is a function that defines all the routes for accessing linked accounts of a certain user
func linkedAccountRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/init.{format:json|xml}", tools.MiddlewareFactory(handler.HandleInitLinkAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/finish.{format:json|xml}", tools.MiddlewareFactory(handler.HandleFinishLinkAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("POST")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetUserLinkedAccounts,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount.{format:json|xml}", tools.MiddlewareFactory(handler.HandleRemoveLinkedAccount,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("DELETE")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/accountinfo/{id}.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetAccountInfo,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/accountprovider.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetAccountProviders,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

	router.HandleFunc("/api/v1/oauth/user/linkedaccount/accountprovider/{id}.{format:json|xml}", tools.MiddlewareFactory(handler.HandleGetAccountProvider,
		handler.Authorization, handler.AuthenticateScope, handler.AccessTokenAuthentication)).Methods("GET")

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

// websocketRoutes is a function that defines all websocket related routes
func websocketRoutes(handler *handler.UserAPIHandler, router *mux.Router) {

	router.HandleFunc("/api/v1/connect.{format:json|xml}/{api_key}/{access_token}", tools.MiddlewareFactory(handler.HandleCreateWebsocket,
		handler.Authorization, handler.WebsocketAccessTokenAuthentication))

	router.HandleFunc("/api/v1/listener/profile", handler.HandleListenToProfileChange).Methods("PUT")

	router.HandleFunc("/api/v1/listener/wallet", handler.HandleListenToWalletChange).Methods("PUT")

	router.HandleFunc("/api/v1/listener/history", handler.HandleListenToHistoryChange).Methods("PUT")

}
