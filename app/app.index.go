package app

import (
	"github.com/Benyam-S/onepay/accountprovider"
	"github.com/Benyam-S/onepay/history"
	"github.com/Benyam-S/onepay/linkedaccount"
	"github.com/Benyam-S/onepay/logger"
	"github.com/Benyam-S/onepay/moneytoken"
	"github.com/Benyam-S/onepay/wallet"
)

// OnePay is a struct that defines all the methods and the fuctions that the onepay system can perform
type OnePay struct {
	WalletService          wallet.IService
	HistoryService         history.IService
	LinkedAccountService   linkedaccount.IService
	MoneyTokenService      moneytoken.IService
	AccountProviderService accountprovider.IService
	Logger                 *logger.Logger
	Channel                chan string
}

// NewApp is a function that creates a new onepay app
func NewApp(walletService wallet.IService, historyService history.IService,
	linkedAccountService linkedaccount.IService, moneyTokenService moneytoken.IService,
	accountProviderService accountprovider.IService, logger *logger.Logger, channel chan string) *OnePay {

	return &OnePay{WalletService: walletService, HistoryService: historyService,
		LinkedAccountService: linkedAccountService, MoneyTokenService: moneyTokenService,
		AccountProviderService: accountProviderService, Logger: logger, Channel: channel}
}
