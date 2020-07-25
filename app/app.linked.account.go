package app

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Benyam-S/onepay/logger"
	"github.com/Benyam-S/onepay/tools"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/middleman"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// GetUserLinkedAccounts is a method that returns a set user's linked accounts
func (onepay *OnePay) GetUserLinkedAccounts(userID string) map[*entity.LinkedAccount]*entity.AccountInfo {

	linkedAccountsMap := make(map[*entity.LinkedAccount]*entity.AccountInfo)
	linkedAccounts := onepay.LinkedAccountService.SearchLinkedAccounts("user_id", userID)

	for _, linkedAccount := range linkedAccounts {
		accountInfo, err := middleman.GetAccountInfo(linkedAccount.AccountID, linkedAccount.AccessToken)
		if err != nil {
			accountInfo = new(entity.AccountInfo)
			accountInfo.AccountID = "error"
			accountInfo.AccountProvider = "error"
		}
		linkedAccountsMap[linkedAccount] = accountInfo
	}

	return linkedAccountsMap

}

// RechargeWallet is a method that recharges user's wallet from external account
func (onepay *OnePay) RechargeWallet(userID, linkedAccountID string, amount float64) error {

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
		return errors.New("insufficient balance, please recharge your wallet")
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

// RemoveLinkedAccount is a method that removes a linked account of a certain user
func (onepay *OnePay) RemoveLinkedAccount(linkedAccountID, userID string) (*entity.LinkedAccount, error) {

	linkedAccount, err := onepay.LinkedAccountService.FindLinkedAccount(linkedAccountID)
	if err != nil {
		return nil, err
	}

	if linkedAccount.UserID != userID {
		return nil, errors.New("linked account doesn't belong to the provided user")
	}

	return onepay.LinkedAccountService.DeleteLinkedAccount(linkedAccountID)

}

// AddLinkedAccount is a method that enables users to add new linked account
func (onepay *OnePay) AddLinkedAccount(userID, accountID, accountProvider string, redisClient *redis.Client) (string, error) {

	linkedAccounts := onepay.LinkedAccountService.SearchLinkedAccounts("account_id", accountID)

	for _, linkedAccount := range linkedAccounts {
		if linkedAccount.AccountProvider == accountProvider {
			return "", errors.New("account has already been linked to other onepay account")
		}
	}

	err := middleman.AddLinkedAccount(accountID, accountProvider)
	if err != nil {
		return "", err
	}

	nonce := uuid.Must(uuid.NewRandom())

	tempOutput, err := json.Marshal(map[string]string{"user_id": userID, "account_id": accountID, "account_provider": accountProvider})
	if err != nil {
		return "", err
	}

	err = tools.SetValue(redisClient, nonce.String(), string(tempOutput), time.Hour*12)
	if err != nil {
		return "", err
	}

	return nonce.String(), nil
}

// VerifyLinkedAccount is a method that verify if the user has inputed a valid nonce with it's otp
func (onepay *OnePay) VerifyLinkedAccount(otp, nonce string, redisClient *redis.Client) error {

	value, err := tools.GetValue(redisClient, nonce)
	if err != nil {
		return errors.New("nonce not found")
	}

	linkedAccountInfo := make(map[string]string)
	err = json.Unmarshal([]byte(value), &linkedAccountInfo)
	if err != nil {
		return err
	}

	// This will be changed
	accessToken, err := middleman.VerifyLinkedAccount(otp)
	if err != nil {
		return err
	}

	newLinkedAccount := new(entity.LinkedAccount)
	newLinkedAccount.AccountID = linkedAccountInfo["account_id"]
	newLinkedAccount.AccountProvider = linkedAccountInfo["account_provider"]
	newLinkedAccount.UserID = linkedAccountInfo["user_id"]
	newLinkedAccount.AccessToken = accessToken

	err = onepay.LinkedAccountService.AddLinkedAccount(newLinkedAccount)
	if err != nil {
		return err
	}

	tools.RemoveValues(redisClient, nonce)

	return nil
}
