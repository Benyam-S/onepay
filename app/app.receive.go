package app

import (
	"errors"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/logger"
)

// ReceiveViaQRCode is a method that endables users to receive money via qr code
func (onepay *OnePay) ReceiveViaQRCode(receiverID string, code string) error {

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

	if moneyToken.SenderID == receiverID {
		return errors.New("cannot make transaction with your own account")
	}

	if moneyToken.Method != entity.MethodTransactionQRCode {
		return errors.New("invalid method, code not found")
	}

	// Delete the money token first in case the use wallet fails to update
	receiverOPWallet.Amount = receiverOPWallet.Amount + moneyToken.Amount
	_, err = onepay.MoneyTokenService.DeleteMoneyToken(moneyToken.Code)
	if err != nil {
		return err
	}

	/* ++++ ++++ +++ checkpoint - wallet +++ ++++ ++++ */
	tempOPWallet := new(entity.UserWallet)
	tempOPWallet.UserID = receiverOPWallet.UserID
	tempOPWallet.Amount = moneyToken.Amount
	logger.Must(onepay.Logger.LogWallet(tempOPWallet))
	/* +++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++++ */

	err = onepay.WalletService.UpdateWallet(receiverOPWallet)
	if err != nil {

		// Adding history for the potential reload
		onepay.AddUserHistory(moneyToken.SenderID, receiverID, entity.MethodTransactionQRCode, moneyToken.Code,
			moneyToken.Amount, moneyToken.SentAt, time.Now())

		return errors.New("wallet checkpoint error")
	}

	/* +++++ +++++ +++++ checkpoint end +++++ ++++ ++++ +++++ */
	logger.Must(onepay.Logger.RemoveWallet(tempOPWallet))
	/* ++++ ++++ ++++ ++++ ++++ ++++ ++++ ++++ +++ ++++ +++++ */

	// Adding history for the received token
	return onepay.AddUserHistory(moneyToken.SenderID, receiverID, entity.MethodTransactionQRCode, moneyToken.Code,
		moneyToken.Amount, moneyToken.SentAt, time.Now())
}
