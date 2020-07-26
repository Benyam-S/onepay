package service

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/linkedaccount"
)

// Service is a type that defines linked account service
type Service struct {
	linkedAccountRepo linkedaccount.ILinkedAccountRepository
}

// NewLinkedAccountService is a function that returns a new linked account service
func NewLinkedAccountService(linkedAccountRepository linkedaccount.ILinkedAccountRepository) linkedaccount.IService {
	return &Service{linkedAccountRepo: linkedAccountRepository}
}

// AddLinkedAccount is a method that adds a new linked account to the system
func (service *Service) AddLinkedAccount(newLinkedAccount *entity.LinkedAccount) error {

	err := service.linkedAccountRepo.Create(newLinkedAccount)
	if err != nil {
		return errors.New("unable to add new linked account")
	}
	return nil
}

// VerifyToLink is a method that verifies if a certain account is ready to linked
// Incase if the database constraint doesn't work
func (service *Service) VerifyToLink(linkedAccount *entity.LinkedAccount) error {

	linkeAccounts := service.linkedAccountRepo.Search("account_id", linkedAccount.AccountID)
	for _, account := range linkeAccounts {
		if account.AccountProvider == linkedAccount.AccountProvider {
			return errors.New("account has been already linked to OnePay user")
		}
	}

	return nil
}

// FindLinkedAccount is a method that find a certain linked account using the provided identifier
func (service *Service) FindLinkedAccount(identifier string) (*entity.LinkedAccount, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	linkedAccount, err := service.linkedAccountRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("linked account not found")
	}

	return linkedAccount, nil
}

// SearchLinkedAccounts is a method that search and returns a set of linked accounts that matchs the identifier value
func (service *Service) SearchLinkedAccounts(columnName string, columnValue interface{}) []*entity.LinkedAccount {
	return service.linkedAccountRepo.Search(columnName, columnValue)
}

// SearchMultipleLinkedAccounts is a method that searchs and returns a set of linked accounts related to the key identifier
func (service *Service) SearchMultipleLinkedAccounts(key, pagination string, columns ...string) []*entity.LinkedAccount {

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return []*entity.LinkedAccount{}
	}

	pageNum, _ := strconv.ParseInt(pagination, 0, 0)
	return service.linkedAccountRepo.SearchMultiple(key, pageNum, columns...)
}

// UpdateLinkedAccount is a method that updates a certain linked account values
func (service *Service) UpdateLinkedAccount(linkedAccount *entity.LinkedAccount) error {

	err := service.linkedAccountRepo.Update(linkedAccount)
	if err != nil {
		return errors.New("unable to update linked account")
	}
	return nil
}

// UpdateLinkedAccountSingleValue is a method that updates a certain linked account's single column value
func (service *Service) UpdateLinkedAccountSingleValue(id, columnName string, columnValue interface{}) error {

	linkedAccount := new(entity.LinkedAccount)
	linkedAccount.ID = id
	err := service.linkedAccountRepo.UpdateValue(linkedAccount, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update linked account")
	}
	return nil
}

// DeleteLinkedAccount is a method that deletes an linked account from the system
func (service *Service) DeleteLinkedAccount(id string) (*entity.LinkedAccount, error) {

	linkedAccount, err := service.linkedAccountRepo.Delete(id)
	if err != nil {
		return nil, errors.New("unable to delete linked account")
	}
	return linkedAccount, nil
}

// DeleteLinkedAccounts is a method that deletes a certain user's linked accounts
func (service *Service) DeleteLinkedAccounts(userID string) ([]*entity.LinkedAccount, error) {

	linkedAccounts, err := service.linkedAccountRepo.DeleteMultiple(userID)
	if err != nil {
		return nil, errors.New("unable to delete linked accounts")
	}
	return linkedAccounts, nil
}
