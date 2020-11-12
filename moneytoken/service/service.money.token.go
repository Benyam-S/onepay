package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/moneytoken"
)

// Service is a type that defines money token service
type Service struct {
	moneyTokenRepo moneytoken.IMoneyTokenRepository
}

// NewMoneyTokenService is a function that returns a new money token service
func NewMoneyTokenService(moneyTokenRepository moneytoken.IMoneyTokenRepository) moneytoken.IService {
	return &Service{moneyTokenRepo: moneyTokenRepository}
}

// AddMoneyToken is a method that adds a new money token to the system
func (service *Service) AddMoneyToken(newMoneyToken *entity.MoneyToken) error {

	// Token will expire after 48 hours
	newMoneyToken.ExpirationDate = time.Now().Add(time.Hour * 48)
	err := service.moneyTokenRepo.Create(newMoneyToken)
	if err != nil {
		return errors.New("unable to add new money token")
	}
	return nil
}

// FindMoneyToken is a method that find a set of money tokens for the provided identifier
func (service *Service) FindMoneyToken(identifier string) (*entity.MoneyToken, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("money token not found")
	}

	moneyToken, err := service.moneyTokenRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("money token not found")
	}
	return moneyToken, nil
}

// SearchMoneyToken is a method that search and returns a set of money tokens for the provided identifier
func (service *Service) SearchMoneyToken(identifier string) []*entity.MoneyToken {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.MoneyToken{}
	}

	return service.moneyTokenRepo.Search(identifier)
}

// ExpiredMoneyTokens is a method that returns all the expired money tokens
func (service *Service) ExpiredMoneyTokens() []*entity.MoneyToken {
	return service.moneyTokenRepo.Expired()
}

// UpdateMoneyToken is a method that updates a certain money token values
func (service *Service) UpdateMoneyToken(moneyToken *entity.MoneyToken) error {

	err := service.moneyTokenRepo.Update(moneyToken)
	if err != nil {
		return errors.New("unable to update money token")
	}
	return nil
}

// UpdateMoneyTokenSingleValue is a method that updates a certain money token's single column value
func (service *Service) UpdateMoneyTokenSingleValue(code, columnName string, columnValue interface{}) error {

	moneyToken := new(entity.MoneyToken)
	moneyToken.Code = code
	err := service.moneyTokenRepo.UpdateValue(moneyToken, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update money token")
	}
	return nil
}

// DeleteMoneyToken is a method that deletes an money token from the system
func (service *Service) DeleteMoneyToken(code string) (*entity.MoneyToken, error) {

	moneyToken, err := service.moneyTokenRepo.Delete(code)
	if err != nil {
		return nil, errors.New("unable to delete money token")
	}
	return moneyToken, nil
}

// DeleteMoneyTokens is a method that deletes a certain user's money tokens
func (service *Service) DeleteMoneyTokens(senderID string) ([]*entity.MoneyToken, error) {

	moneyTokens, err := service.moneyTokenRepo.DeleteMultiple(senderID)
	if err != nil {
		return nil, errors.New("unable to delete money tokens")
	}
	return moneyTokens, nil
}
