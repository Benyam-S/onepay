package app

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/logger"
)

// ReloadHistory is a method that enables the server to reload any history adding process
func (onepay *OnePay) ReloadHistory() {

	loggedHistories := onepay.Logger.LoggedHistories()
	for _, loggedHistory := range loggedHistories {

		tempHistory := new(entity.UserHistory)
		tempHistory.Amount = loggedHistory.Amount
		tempHistory.Code = loggedHistory.Code
		tempHistory.Method = loggedHistory.Method
		tempHistory.SentAt = loggedHistory.SentAt
		tempHistory.ReceivedAt = loggedHistory.ReceivedAt
		tempHistory.SenderID = loggedHistory.SenderID
		tempHistory.ReceiverID = loggedHistory.ReceiverID

		err := onepay.HistoryService.AddHistory(loggedHistory)
		if err == nil {
			logger.Must(onepay.Logger.RemoveHistory(tempHistory))
		}
	}
}

// ReloadMoneyToken is a method that enables the server to reload any money token adding process
func (onepay *OnePay) ReloadMoneyToken() {

	loggedMoneyTokens := onepay.Logger.LoggedMoneyTokens()
	for _, loggedMoneyToken := range loggedMoneyTokens {

		tempMoneyToken := new(entity.MoneyToken)
		tempMoneyToken.Amount = loggedMoneyToken.Amount
		tempMoneyToken.Code = loggedMoneyToken.Code
		tempMoneyToken.ExpirationDate = loggedMoneyToken.ExpirationDate
		tempMoneyToken.Method = loggedMoneyToken.Method
		tempMoneyToken.SenderID = loggedMoneyToken.SenderID
		tempMoneyToken.SentAt = loggedMoneyToken.SentAt

		err := onepay.MoneyTokenService.AddMoneyToken(loggedMoneyToken)
		if err == nil {
			logger.Must(onepay.Logger.RemoveMoneyToken(tempMoneyToken))
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

		// This will protect the user from multiple draining
		if opWallet.Amount <= 0 && loggedWallet.Amount < 0 {
			logger.Must(onepay.Logger.RemoveWallet(loggedWallet))
			continue
		}

		opWallet.Amount += loggedWallet.Amount
		err = onepay.WalletService.UpdateWallet(opWallet)
		if err == nil {
			logger.Must(onepay.Logger.RemoveWallet(loggedWallet))
		}
	}
}
