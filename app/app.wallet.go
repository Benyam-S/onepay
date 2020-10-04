package app

import (
	"errors"
	"time"

	"github.com/Benyam-S/onepay/middleman"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/logger"
)

// DrainWallet is a method that drains all the cash out your wallet
func (onepay *OnePay) DrainWallet(userID, linkedAccountID string) error {

	opWallet, err := onepay.WalletService.FindWallet(userID)
	if err != nil {
		return err
	}

	amount := opWallet.Amount
	if amount <= 0 {
		return errors.New("can not drain empty wallet")
	}
	// draning the account
	opWallet.Amount = 0

	linkedAccount, err := onepay.LinkedAccountService.FindLinkedAccount(linkedAccountID)
	if err != nil {
		return err
	}

	if linkedAccount.UserID != userID {
		return errors.New("linked account doesn't belong to the provided user")
	}

	err = onepay.WalletService.UpdateWallet(opWallet)
	if err != nil {
		return err
	}

	/* ++++ ++++ +++ checkpoint - wallet +++ ++++ ++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = opWallet.UserID
	tempOPWallet.Amount = amount
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++++ */

	err = middleman.RefillAccount(linkedAccount.AccountID, linkedAccount.AccessToken, amount)
	if err != nil {

		/* +++++++++++++++++++++++ Undo +++++++++++++++++++++++ */
		opWallet.Amount += amount
		innerErr := onepay.WalletService.UpdateWallet(opWallet)
		/* ++++++++++++++++++++++++++++++++++++++++++++++++++++ */

		if innerErr != nil {

			// Adding history for the potential reload
			onepay.AddUserHistory(userID, linkedAccount.AccountID, entity.MethodWithdrawn, "",
				amount, time.Now(), time.Now())

			return errors.New(entity.WalletCheckpointError)
		}

		/* +++++ +++++ +++++ checkpoint end +++++ +++++ +++++ */
		logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
		/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

		return err
	}

	/* +++++ +++++ +++++ checkpoint end +++++ +++++ +++++ */
	logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	// Adding history for the withdrawal process
	return onepay.AddUserHistory(userID, linkedAccount.AccountID, entity.MethodWithdrawn, "",
		amount, time.Now(), time.Now())
}

// MarkWalletAsViewed is a method that marks the wallet change as viewed
func (onepay *OnePay) MarkWalletAsViewed(userID string) error {

	err := onepay.WalletService.UpdateWalletSeen(userID, true)
	if err != nil {
		return err
	}

	return nil
}

// AddUserHistory is a method that add user history for the onepay app methods
func (onepay *OnePay) AddUserHistory(senderID, receiverID, method, code string,
	amount float64, sentAt, receivedAt time.Time) error {

	opHistory := new(entity.UserHistory)
	opHistory.Amount = amount
	opHistory.Code = code
	opHistory.Method = method
	opHistory.SentAt = sentAt
	opHistory.ReceivedAt = receivedAt
	opHistory.SenderID = senderID
	opHistory.ReceiverID = receiverID

	/* +++++ +++++ checkpoint - history +++++ ++++++ */
	// tempHistory is created because the .AddHistory() method wil change some value's of the opHistory object
	tempHistory := new(entity.UserHistory)
	tempHistory.Amount = amount
	tempHistory.Code = code
	tempHistory.Method = method
	tempHistory.SentAt = sentAt
	tempHistory.ReceivedAt = receivedAt
	tempHistory.SenderID = senderID
	tempHistory.ReceiverID = receiverID
	logger.Must(onepay.Logger.LogHistory(tempHistory))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	err := onepay.HistoryService.AddHistory(opHistory)
	if err != nil {
		return errors.New(entity.HistoryCheckpointError)
	}

	/* +++++ +++++ +++++ checkpoint end ++++ ++++ +++++ */
	logger.Must(onepay.Logger.RemoveHistory(tempHistory))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++ */

	return nil
}

// RechargeWallet is a method that recharges user's wallet from external account
func (onepay *OnePay) RechargeWallet(userID, linkedAccountID string, amount float64) error {

	opWallet, err := onepay.WalletService.FindWallet(userID)
	if err != nil {
		return errors.New("user wallet not found")
	}

	linkedAccount, err := onepay.LinkedAccountService.FindLinkedAccount(linkedAccountID)
	if err != nil {
		return errors.New("linked account not found")
	}

	if linkedAccount.UserID != userID {
		// It is appropirate to this that "linked account doesn't belong to the provided user"
		return errors.New("linked account not found")
	}

	accountInfo, err := middleman.GetAccountInfo(linkedAccount.AccountID, linkedAccount.AccessToken)
	if err != nil {
		return err
	}

	if accountInfo.Amount < amount {
		return errors.New("insufficient balance, please recharge your linked account")
	}

	err = middleman.WithdrawFromAccount(linkedAccount.AccountID, linkedAccount.AccessToken, amount)
	if err != nil {
		return err
	}

	opWallet.Amount += amount

	/* ++++ ++++ +++ checkpoint - wallet +++ ++++ ++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = opWallet.UserID
	tempOPWallet.Amount = amount
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ +++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	err = onepay.WalletService.UpdateWallet(opWallet)
	if err != nil {

		// Adding history for the potential reload
		onepay.AddUserHistory(linkedAccount.AccountID, userID, entity.MethodRecharged, "",
			amount, time.Now(), time.Now())

		return errors.New(entity.WalletCheckpointError)
	}

	/* +++++ +++++ +++++ checkpoint end +++++ +++++ +++++ */
	logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	// Adding history for the recharging process
	return onepay.AddUserHistory(linkedAccount.AccountID, userID, entity.MethodRecharged, "",
		amount, time.Now(), time.Now())
}

// WithdrawFromWallet is a method that enables user's to withdraw money from onepay account/wallet
func (onepay *OnePay) WithdrawFromWallet(userID, linkedAccountID string, amount float64) error {

	opWallet, err := onepay.WalletService.FindWallet(userID)
	if err != nil {
		return err
	}

	linkedAccount, err := onepay.LinkedAccountService.FindLinkedAccount(linkedAccountID)
	if err != nil {
		return err
	}

	if linkedAccount.UserID != userID {
		return errors.New("linked account doesn't belong to the provided user")
	}

	if !AboveWithdrawBaseLimit(amount) {
		return errors.New("the provided amount is less than the withdraw base limit")
	}

	if opWallet.Amount < amount {
		return errors.New(entity.InsufficientBalanceError)
	}

	opWallet.Amount -= amount

	err = onepay.WalletService.UpdateWallet(opWallet)
	if err != nil {
		return err
	}

	/* ++++ ++++ +++ checkpoint - wallet +++ ++++ ++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = opWallet.UserID
	tempOPWallet.Amount = amount
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ +++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	err = middleman.RefillAccount(linkedAccount.AccountID, linkedAccount.AccessToken, amount)
	if err != nil {

		/* +++++++++++++++++++++++ Undo +++++++++++++++++++++++ */
		opWallet.Amount += amount
		innerErr := onepay.WalletService.UpdateWallet(opWallet)
		/* ++++++++++++++++++++++++++++++++++++++++++++++++++++ */

		if innerErr != nil {

			// Adding history for the potential reload
			onepay.AddUserHistory(userID, linkedAccount.AccountID, entity.MethodWithdrawn, "",
				amount, time.Now(), time.Now())

			return errors.New(entity.WalletCheckpointError)
		}

		/* +++++ +++++ +++++ checkpoint end +++++ +++++ +++++ */
		logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
		/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

		return err
	}

	/* +++++ +++++ +++++ checkpoint end +++++ +++++ +++++ */
	logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	// Adding history for the withdrawal process
	return onepay.AddUserHistory(userID, linkedAccount.AccountID, entity.MethodWithdrawn, "",
		amount, time.Now(), time.Now())
}
