package app

import (
	"errors"
	"time"

	"github.com/Benyam-S/onepay/middleman"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/logger"
)

// UserHistory is a method that returns the user history only
func (onepay *OnePay) UserHistory(userID string, pagenation int64, viewBys ...string) []*entity.UserHistory {

	orderBy := "id"
	opHistories := make([]*entity.UserHistory, 0)
	length := len(viewBys)

	for _, viewBy := range viewBys {

		searchColumns := make([]string, 0)
		methods := make([]string, 0)

		if viewBy == "transaction_received" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode)

		} else if viewBy == "transaction_sent" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode)

		} else if viewBy == "payment_received" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodPaymentQRCode)

		} else if viewBy == "payment_sent" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodPaymentQRCode)

		} else if viewBy == "recharged" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodRecharged)

		} else if viewBy == "withdrawn" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodWithdrawn)

		} else if viewBy == "all" && length == 1 {
			searchColumns = append(searchColumns, "sender_id", "receiver_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode, entity.MethodPaymentQRCode)
		} else {
			// If it is unkown view by
			continue
		}

		result := onepay.HistoryService.SearchHistories(userID, orderBy, methods, pagenation, searchColumns...)
		opHistories = append(opHistories, result...)

	}

	return opHistories
}

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
