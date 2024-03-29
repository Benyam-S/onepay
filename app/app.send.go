package app

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/logger"
	"github.com/go-redis/redis"

	"github.com/Benyam-S/onepay/entity"
)

// SendViaQRCode is a method that enables user to send money via qr code
func (onepay *OnePay) SendViaQRCode(userID string, amount float64, redisClient *redis.Client) (*entity.MoneyToken, error) {

	if !AboveTransactionBaseLimit(amount) {
		return nil, errors.New(entity.TransactionBaseLimitError)
	}

	if AboveDailyTransactionLimit(userID, amount, redisClient) {
		return nil, errors.New(entity.DailyTransactionLimitError)
	}

	opWallet, err := onepay.WalletService.FindWallet(userID)
	if err != nil {
		return nil, err
	}

	transactionFee := GetTransactionFee(amount)
	if opWallet.Amount < amount+transactionFee {
		return nil, errors.New(entity.InsufficientBalanceError)
	}

	opWallet.Amount = opWallet.Amount - (amount + transactionFee)
	err = onepay.WalletService.UpdateWallet(opWallet)
	if err != nil {
		return nil, err
	}

	moneyToken := new(entity.MoneyToken)
	moneyToken.Amount = amount
	moneyToken.Method = entity.MethodTransactionQRCode
	moneyToken.SenderID = opWallet.UserID
	moneyToken.SentAt = time.Now()

	/* ++++ ++++ +++ checkpoint - money token ++++ ++++ +++ */
	tempMoneyToken := new(entity.MoneyToken)
	tempMoneyToken.Amount = amount
	tempMoneyToken.Method = entity.MethodTransactionQRCode
	tempMoneyToken.SenderID = opWallet.UserID
	tempMoneyToken.SentAt = time.Now()
	logger.Must(onepay.Logger.LogMoneyToken(tempMoneyToken))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++ */

	err = onepay.MoneyTokenService.AddMoneyToken(moneyToken)
	if err != nil {

		/* ++++++++++++++++++++++++++ Undo ++++++++++++++++++++++++++ */
		opWallet.Amount = opWallet.Amount + (amount + transactionFee)
		innerErr := onepay.WalletService.UpdateWallet(opWallet)
		/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

		if innerErr != nil {
			return nil, errors.New(entity.MoneyTokenCheckpointError)
		}

		/* +++++ ++++ ++++ ++++ checkpoint end ++++ ++++ ++++ +++++ */
		logger.Must(onepay.Logger.RemoveMoneyToken(tempMoneyToken))
		/* +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */
		return nil, err
	}

	/* +++++ ++++ ++++ ++++ checkpoint end ++++ ++++ ++++ +++++ */
	logger.Must(onepay.Logger.RemoveMoneyToken(tempMoneyToken))
	/* +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	// Just updating the users daily transaction limit
	AddToDailyTransaction(userID, amount, redisClient)

	return moneyToken, nil

}

// SendViaOnePayID is a method that enables user to send money via onepay id
func (onepay *OnePay) SendViaOnePayID(senderID, receiverID string,
	amount float64, redisClient *redis.Client) error {

	if !AboveTransactionBaseLimit(amount) {
		return errors.New(entity.TransactionBaseLimitError)
	}

	if AboveDailyTransactionLimit(senderID, amount, redisClient) {
		return errors.New(entity.DailyTransactionLimitError)
	}

	if senderID == receiverID {
		return errors.New(entity.TransactionWSelfError)
	}

	senderOPWallet, err := onepay.WalletService.FindWallet(senderID)
	if err != nil {
		return errors.New(entity.SenderNotFoundError)
	}

	receiverOPWallet, err := onepay.WalletService.FindWallet(receiverID)
	if err != nil {
		return errors.New(entity.ReceiverNotFoundError)
	}

	transactionFee, _ := strconv.ParseFloat(os.Getenv(entity.TransactionFee), 64)
	if senderOPWallet.Amount < amount+transactionFee {
		return errors.New(entity.InsufficientBalanceError)
	}

	senderOPWallet.Amount = senderOPWallet.Amount - (amount + transactionFee)
	receiverOPWallet.Amount = receiverOPWallet.Amount + amount

	err = onepay.WalletService.UpdateWallet(senderOPWallet)
	if err != nil {
		return err
	}

	/* ++++ ++++ +++ checkpoint - wallet +++ ++++ ++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = receiverOPWallet.UserID
	tempOPWallet.Amount = amount
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++++ */

	err = onepay.WalletService.UpdateWallet(receiverOPWallet)
	if err != nil {

		/* ++++++++++++++++++++++++++++++++ Undo +++++++++++++++++++++++++++++++ */
		senderOPWallet.Amount = senderOPWallet.Amount + (amount + transactionFee)
		innerErr := onepay.WalletService.UpdateWallet(senderOPWallet)
		/* +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

		if innerErr != nil {

			// Adding history for the potential reload
			onepay.AddUserHistory(senderID, receiverID, entity.MethodTransactionOnePayID, "",
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

	// Just updating the users daily transaction limit
	AddToDailyTransaction(senderID, amount, redisClient)

	// Adding history for the given transaction
	return onepay.AddUserHistory(senderID, receiverID, entity.MethodTransactionOnePayID, "",
		amount, time.Now(), time.Now())

}
