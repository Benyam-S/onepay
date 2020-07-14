package service

import (
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
	return service.moneyTokenRepo.Create(newMoneyToken)
}

// FindMoneyToken is a method that find a set of money tokens for the provided identifier
func (service *Service) FindMoneyToken(identifier string) (*entity.MoneyToken, error) {
	return service.moneyTokenRepo.Find(identifier)
}

// SearchMoneyToken is a method that search and returns a set of money tokens for the provided identifier
func (service *Service) SearchMoneyToken(identifier string) []*entity.MoneyToken {
	return service.moneyTokenRepo.Search(identifier)
}

// UpdateMoneyToken is a method that updates a certain money token values
func (service *Service) UpdateMoneyToken(moneyToken *entity.MoneyToken) error {
	return service.moneyTokenRepo.Update(moneyToken)
}

// UpdateMoneyTokenSingleValue is a method that updates a certain money token's single column value
func (service *Service) UpdateMoneyTokenSingleValue(code, columnName string, columnValue interface{}) error {
	moneyToken := new(entity.MoneyToken)
	moneyToken.Code = code
	return service.moneyTokenRepo.UpdateValue(moneyToken, columnName, columnValue)
}

// DeleteMoneyToken is a method that deletes an money token from the system
func (service *Service) DeleteMoneyToken(code string) (*entity.MoneyToken, error) {

	return service.moneyTokenRepo.Delete(code)
}

// DeleteMoneyTokens is a method that deletes a certain user's money tokens
func (service *Service) DeleteMoneyTokens(senderID string) ([]*entity.MoneyToken, error) {

	return service.moneyTokenRepo.DeleteMultiple(senderID)
}
