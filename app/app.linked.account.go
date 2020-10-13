package app

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Benyam-S/onepay/tools"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/middleman"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// GetUserLinkedAccounts is a method that returns a set user's linked accounts
func (onepay *OnePay) GetUserLinkedAccounts(userID string) []*entity.LinkedAccountContainer {

	linkedAccounts := onepay.LinkedAccountService.SearchLinkedAccounts("user_id", userID)
	linkedAccountsFiltered := make([]*entity.LinkedAccountContainer, 0)

	for _, linkedAccount := range linkedAccounts {

		accountProvider, err := onepay.AccountProviderService.FindAccountProvider(linkedAccount.AccountProviderID)
		if err != nil {
			continue
		}

		filteredLinkedAccount := new(entity.LinkedAccountContainer)
		filteredLinkedAccount.ID = linkedAccount.ID
		filteredLinkedAccount.UserID = linkedAccount.UserID
		filteredLinkedAccount.AccountID = linkedAccount.AccountID
		filteredLinkedAccount.AccountProviderID = linkedAccount.AccountProviderID
		filteredLinkedAccount.AccountProviderName = accountProvider.Name

		linkedAccountsFiltered = append(linkedAccountsFiltered, filteredLinkedAccount)
	}

	return linkedAccountsFiltered

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
func (onepay *OnePay) AddLinkedAccount(userID, accountID, accountProviderID string, redisClient *redis.Client) (string, error) {

	linkedAccounts := onepay.LinkedAccountService.SearchLinkedAccounts("account_id", accountID)

	for _, linkedAccount := range linkedAccounts {
		if linkedAccount.AccountProviderID == accountProviderID {
			return "", errors.New("account has already been linked to other onepay account")
		}
	}

	// we can insert this method to middleman package if needed
	// but for now all we need to check is the availablity of the account provider
	_, err := onepay.AccountProviderService.FindAccountProvider(accountProviderID)
	if err != nil {
		return "", err
	}

	err = middleman.AddLinkedAccount(accountID, accountProviderID)
	if err != nil {
		return "", err
	}

	nonce := uuid.Must(uuid.NewRandom())

	tempOutput, err := json.Marshal(map[string]string{"user_id": userID,
		"account_id": accountID, "account_provider_id": accountProviderID})
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
func (onepay *OnePay) VerifyLinkedAccount(otp, nonce string, redisClient *redis.Client) (*entity.LinkedAccountContainer, error) {

	value, err := tools.GetValue(redisClient, nonce)
	if err != nil {
		return nil, errors.New("nonce not found")
	}

	linkedAccountInfo := make(map[string]string)
	err = json.Unmarshal([]byte(value), &linkedAccountInfo)
	if err != nil {
		return nil, err
	}

	// This will be changed
	accessToken, err := middleman.VerifyLinkedAccount(otp)
	if err != nil {
		return nil, err
	}

	// we can insert this method to middleman package if needed
	// but for now all we need to check is the availablity of the account provider
	_, err = onepay.AccountProviderService.FindAccountProvider(linkedAccountInfo["account_provider_id"])
	if err != nil {
		return nil, err
	}

	newLinkedAccount := new(entity.LinkedAccount)
	newLinkedAccount.AccountID = linkedAccountInfo["account_id"]
	newLinkedAccount.AccountProviderID = linkedAccountInfo["account_provider_id"]
	newLinkedAccount.UserID = linkedAccountInfo["user_id"]
	newLinkedAccount.AccessToken = accessToken

	err = onepay.LinkedAccountService.AddLinkedAccount(newLinkedAccount)
	if err != nil {
		return nil, err
	}

	tools.RemoveValues(redisClient, nonce)

	accountProvider, _ := onepay.AccountProviderService.FindAccountProvider(newLinkedAccount.AccountProviderID)
	linkedAccount := new(entity.LinkedAccountContainer)
	linkedAccount.ID = newLinkedAccount.ID
	linkedAccount.UserID = newLinkedAccount.UserID
	linkedAccount.AccountID = newLinkedAccount.AccountID
	linkedAccount.AccountProviderID = newLinkedAccount.AccountProviderID
	linkedAccount.AccountProviderName = accountProvider.Name

	return linkedAccount, nil
}
