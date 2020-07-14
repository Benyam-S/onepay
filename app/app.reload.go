package app

import (
	"time"

	"github.com/Benyam-S/onepay/logger"
)

// ReloadHistory is a method that enables the server to reload any history adding process
func (onepay *OnePay) ReloadHistory() {

	loggedHistories := onepay.Logger.LoggedHistories()
	for _, loggedHistory := range loggedHistories {
		err := onepay.HistoryService.AddHistory(loggedHistory)
		if err != nil {
			logger.Must(onepay.Logger.RemoveHistory(loggedHistory))
		}
	}
}

// ReloadMoneyToken is a method that enables the server to reload any money token adding process
func (onepay *OnePay) ReloadMoneyToken() {

	loggedMoneyTokens := onepay.Logger.LoggedMoneyTokens()
	for _, loggedMoneyToken := range loggedMoneyTokens {
		err := onepay.MoneyTokenService.AddMoneyToken(loggedMoneyToken)
		if err == nil {
			logger.Must(onepay.Logger.RemoveMoneyToken(loggedMoneyToken))
		}
	}
}

// ReloadWallet is a method that enables the server to reload any user wallet updating process
func (onepay *OnePay) ReloadWallet() {

	loggedWallets := onepay.Logger.LoggedWallets()
	for _, loggedWallet := range loggedWallets {
		opWallet, err := onepay.WalletService.FindWallet(loggedWallet.UserID)
		if err != nil {
			continue
		}

		opWallet.Amount += loggedWallet.Amount
		err = onepay.WalletService.UpdateWallet(opWallet)
		if err == nil {
			logger.Must(onepay.Logger.RemoveWallet(opWallet))
		}
	}
}

// Reload is a reload method that encompass all the reload methods
func (onepay *OnePay) Reload() {

	go func() {
		for {
			time.Sleep(time.Hour * 2)
			onepay.Channel <- "all"
		}
	}()

	for {

		value := <-onepay.Channel
		switch value {

		case "all":
			onepay.ReloadMoneyToken()
			onepay.ReloadWallet()
			onepay.ReloadHistory()

		case "reload_money_token":
			onepay.ReloadMoneyToken()

		case "reload_wallet":
			onepay.ReloadWallet()
			fallthrough

		case "reload_history":
			onepay.ReloadHistory()
		}
	}
}
