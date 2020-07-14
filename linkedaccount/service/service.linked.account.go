package service

import (
	"errors"

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
	return service.linkedAccountRepo.Create(newLinkedAccount)
}

// VerifyToLink is a method that verifies if a certain account is ready to linked
// Incase if the database constraint doesn't work
func (service *Service) VerifyToLink(linkedAccount *entity.LinkedAccount) error {
	linkeAccounts := service.linkedAccountRepo.Search("account_id", linkedAccount.AccountID)
	for _, account := range linkeAccounts {
		if account.AccountProvider == linkedAccount.AccountProvider {
			return errors.New("account has been already linked to other OnePay user")
		}
	}

	return nil
}

// FindLinkedAccount is a method that find a certain linked account using the provided identifier
func (service *Service) FindLinkedAccount(identifier string) (*entity.LinkedAccount, error) {
	return service.linkedAccountRepo.Find(identifier)
}

// SearchLinkedAccounts is a method that search and returns a set of linked accounts that matchs the identifier value
func (service *Service) SearchLinkedAccounts(columnName string, columnValue interface{}) []*entity.LinkedAccount {
	return service.linkedAccountRepo.Search(columnName, columnValue)
}

// UpdateLinkedAccount is a method that updates a certain linked account values
func (service *Service) UpdateLinkedAccount(linkedAccount *entity.LinkedAccount) error {
	return service.linkedAccountRepo.Update(linkedAccount)
}

// UpdateLinkedAccountSingleValue is a method that updates a certain linked account's single column value
func (service *Service) UpdateLinkedAccountSingleValue(id, columnName string, columnValue interface{}) error {
	linkedAccount := new(entity.LinkedAccount)
	linkedAccount.ID = id
	return service.linkedAccountRepo.UpdateValue(linkedAccount, columnName, columnValue)
}

// DeleteLinkedAccount is a method that deletes an linked account from the system
func (service *Service) DeleteLinkedAccount(id string) (*entity.LinkedAccount, error) {
	return service.linkedAccountRepo.Delete(id)
}

// DeleteLinkedAccounts is a method that deletes a certain user's linked accounts
func (service *Service) DeleteLinkedAccounts(userID string) ([]*entity.LinkedAccount, error) {
	return service.linkedAccountRepo.DeleteMultiple(userID)
}
