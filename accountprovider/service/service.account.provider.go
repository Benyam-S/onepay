package service

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/onepay/accountprovider"
	"github.com/Benyam-S/onepay/entity"
)

// Service is a type that defines account provider service
type Service struct {
	accountProviderRepo accountprovider.IAccountProviderRepository
}

// NewAccountProviderService is a function that returns a new account provider service
func NewAccountProviderService(
	accountProviderRepository accountprovider.IAccountProviderRepository) accountprovider.IService {
	return &Service{accountProviderRepo: accountProviderRepository}
}

// AddAccountProvider is a method that adds a new account provider to the system
func (service *Service) AddAccountProvider(newAccountProvider *entity.AccountProvider) error {

	err := service.accountProviderRepo.Create(newAccountProvider)
	if err != nil {
		return errors.New("unable to add new account provider")
	}
	return nil
}

// ValidateAccountProvider is a method that validates a certain account provider entries
func (service *Service) ValidateAccountProvider(accountProvider *entity.AccountProvider) entity.ErrMap {

	errMap := make(map[string]error)

	emptyAccountProviderName, _ := regexp.MatchString(`^\s*$`, accountProvider.Name)
	if emptyAccountProviderName {
		errMap["name"] = errors.New("invalid account provider name used")
	}

	// Meaning new account provider
	if accountProvider.ID == "" {
		if !emptyAccountProviderName && !service.accountProviderRepo.IsUnique("name", accountProvider.Name) {
			errMap["name"] = errors.New("account provider with the provided name already exists")
		}
	} else {
		tempAccountProvider, _ := service.accountProviderRepo.Find(accountProvider.ID)
		if tempAccountProvider.Name != accountProvider.Name {
			if !emptyAccountProviderName &&
				!service.accountProviderRepo.IsUnique("name", accountProvider.Name) {
				errMap["name"] = errors.New("account provider with the provided name already exists")
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindAccountProvider is a method that find an account provider for the provided identifier
func (service *Service) FindAccountProvider(identifier string) (*entity.AccountProvider, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	accountProvider, err := service.accountProviderRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("account provider not found")
	}
	return accountProvider, nil
}

// AllAccountProviders is a method that returns all the account providers registered in the system
func (service *Service) AllAccountProviders() []*entity.AccountProvider {
	return service.accountProviderRepo.All()
}

// SearchAccountProvider is a method that search and returns a set of account providers for the provided identifier
func (service *Service) SearchAccountProvider(key, orderBy string, pageNum int64, columns ...string) ([]*entity.AccountProvider, int64) {

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return []*entity.AccountProvider{}, 0
	}

	return service.accountProviderRepo.SearchMultiple(key, orderBy, pageNum, columns...)
}

// UpdateAccountProvider is a method that updates a certain account provider values
func (service *Service) UpdateAccountProvider(accountProvider *entity.AccountProvider) error {

	err := service.accountProviderRepo.Update(accountProvider)
	if err != nil {
		return errors.New("unable to update account provider")
	}
	return nil
}

// UpdateAccountProviderSingleValue is a method that updates a certain account provider's single column value
func (service *Service) UpdateAccountProviderSingleValue(id, columnName string, columnValue interface{}) error {

	accountProvider := new(entity.AccountProvider)
	accountProvider.ID = id

	err := service.accountProviderRepo.UpdateValue(accountProvider, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update account provider")
	}
	return nil
}

// DeleteAccountProvider is a method that deletes an account provider from the system
func (service *Service) DeleteAccountProvider(id string) (*entity.AccountProvider, error) {

	accountProvider, err := service.accountProviderRepo.Delete(id)
	if err != nil {
		return nil, errors.New("unable to delete account provider")
	}
	return accountProvider, nil
}
