package app

import (
	"errors"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/logger"
	"github.com/go-redis/redis"
)

// CreatePaymentToken is a method that generate a payment token
func (onepay *OnePay) CreatePaymentToken(userID string, amount float64) (*entity.MoneyToken, error) {

	if !AboveTransactionBaseLimit(amount) {
		return nil, errors.New("the provided amount is less than the transaction base limit")
	}

	opWallet, err := onepay.WalletService.FindWallet(userID)
	if err != nil {
		return nil, err
	}

	moneyToken := new(entity.MoneyToken)
	moneyToken.Amount = amount
	moneyToken.Method = entity.MethodPaymentQRCode
	moneyToken.SenderID = opWallet.UserID
	moneyToken.SentAt = time.Now()

	err = onepay.MoneyTokenService.AddMoneyToken(moneyToken)
	if err != nil {
		return nil, err
	}

	return moneyToken, nil
}

// PayViaQRCode is a method that enables user's to pay to another user via qr code.
// InPayViaQRCode the receiverID is referred to as the one who is paying the money, in the other word receiving  the qr code.
// Also transaction fee while be deducted from the senderID since the sender initiated the request
func (onepay *OnePay) PayViaQRCode(receiverID string, code string, redisClient *redis.Client) error {

	receiverOPWallet, err := onepay.WalletService.FindWallet(receiverID)
	if err != nil {
		return errors.New("unable to find user for the provided receiver id")
	}

	moneyToken, err := onepay.MoneyTokenService.FindMoneyToken(code)
	if err != nil {
		return err
	}

	if !moneyToken.ExpirationDate.After(time.Now()) {
		return errors.New("money token had past expiration date")
	}

	if !AboveTransactionBaseLimit(moneyToken.Amount) {
		return errors.New("the provided amount is less than the transaction base limit")
	}

	if AboveDailyTransactionLimit(receiverID, moneyToken.Amount, redisClient) {
		return errors.New("user has exceeded daily transaction limit")
	}

	if moneyToken.SenderID == receiverID {
		return errors.New("cannot make transaction with your own account")
	}

	if moneyToken.Method != entity.MethodPaymentQRCode {
		return errors.New("invalid method, code not found")
	}

	senderOPWallet, err := onepay.WalletService.FindWallet(moneyToken.SenderID)
	if err != nil {
		return errors.New("unable to find user for the provided sender id")
	}

	if receiverOPWallet.Amount < moneyToken.Amount {
		return errors.New("insufficient balance, please recharge your wallet")
	}

	transactionFee := GetTransactionFee(moneyToken.Amount)
	receiverOPWallet.Amount = receiverOPWallet.Amount - moneyToken.Amount
	senderOPWallet.Amount = senderOPWallet.Amount + (moneyToken.Amount - transactionFee)

	// Delete the money token first
	_, err = onepay.MoneyTokenService.DeleteMoneyToken(moneyToken.Code)
	if err != nil {
		return err
	}

	err = onepay.WalletService.UpdateWallet(receiverOPWallet)
	if err != nil {
		return err
	}

	/* +++++ +++++ +++++ checkpoint - wallet ++++ ++++ +++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = senderOPWallet.UserID
	tempOPWallet.Amount = moneyToken.Amount - transactionFee
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ +++++ +++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	err = onepay.WalletService.UpdateWallet(senderOPWallet)
	if err != nil {

		/* ++++++++++++++++++++++++++++++ Undo ++++++++++++++++++++++++++++++ */
		receiverOPWallet.Amount = receiverOPWallet.Amount + moneyToken.Amount
		innerErr := onepay.WalletService.UpdateWallet(receiverOPWallet)
		/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

		if innerErr != nil {

			// Adding history for the potential reload
			onepay.AddUserHistory(receiverID, moneyToken.SenderID, entity.MethodPaymentQRCode, moneyToken.Code,
				moneyToken.Amount, moneyToken.SentAt, time.Now())

			return errors.New("wallet checkpoint error")
		}

		/* ++++ ++++ ++++ +++ checkpoint end +++ ++++ ++++ ++++ */
		logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
		/* +++++ +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

		return err
	}

	/* ++++ ++++ ++++ +++ checkpoint end +++ ++++ ++++ ++++ */
	logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
	/* +++++ +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++++ */

	// Just updating the users daily transaction limit
	AddToDailyTransaction(receiverID, moneyToken.Amount, redisClient)

	// Adding history for the received payment
	return onepay.AddUserHistory(receiverID, moneyToken.SenderID, entity.MethodPaymentQRCode, moneyToken.Code,
		moneyToken.Amount, moneyToken.SentAt, time.Now())

}
