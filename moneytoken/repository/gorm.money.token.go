package repository

import (
	"errors"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/moneytoken"
	"github.com/Benyam-S/onepay/tools"
	"github.com/jinzhu/gorm"
)

// MoneyTokenRepository is a type that defines a money token repository
type MoneyTokenRepository struct {
	conn *gorm.DB
}

// NewMoneyTokenRepository is a function that returns a new money token repository
func NewMoneyTokenRepository(connection *gorm.DB) moneytoken.IMoneyTokenRepository {
	return &MoneyTokenRepository{conn: connection}
}

// Create is a method that adds a new money token to the database
func (repo *MoneyTokenRepository) Create(newMoneyToken *entity.MoneyToken) error {

	newMoneyToken.Code = tools.GenerateMoneyTokenCode()

	for !repo.IsUnique("code", newMoneyToken.Code) {
		newMoneyToken.Code = tools.GenerateMoneyTokenCode()
	}

	err := repo.conn.Create(newMoneyToken).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain money token from the database using an identifier.
// In Find() code is only used as a key
func (repo *MoneyTokenRepository) Find(identifier string) (*entity.MoneyToken, error) {
	moneyToken := new(entity.MoneyToken)
	err := repo.conn.Model(moneyToken).
		Where("code = ?", identifier).
		First(moneyToken).Error

	if err != nil {
		return nil, err
	}
	return moneyToken, nil
}

// Search is a method that search and returns a set of money tokens that is limited to the provided identifier
// Search uses sender_id only as a key since there is no other important entity left used for searching
func (repo *MoneyTokenRepository) Search(identifier string) []*entity.MoneyToken {
	var moneyTokens []*entity.MoneyToken
	err := repo.conn.Model(entity.MoneyToken{}).
		Where("sender_id = ?", identifier).
		Find(&moneyTokens).Error

	if err != nil {
		return []*entity.MoneyToken{}
	}
	return moneyTokens
}

// Expired is a method that returns all the expired moneytokens
func (repo *MoneyTokenRepository) Expired() []*entity.MoneyToken {
	var moneyTokens []*entity.MoneyToken
	err := repo.conn.Model(entity.MoneyToken{}).
		Where("expiration_date < ?", time.Now()).
		Find(&moneyTokens).Error

	if err != nil {
		return []*entity.MoneyToken{}
	}
	return moneyTokens
}

// Update is a method that updates a certain money token value in the database
func (repo *MoneyTokenRepository) Update(moneyToken *entity.MoneyToken) error {

	prevMoneyToken := new(entity.MoneyToken)
	err := repo.conn.Model(prevMoneyToken).Where("code = ?", moneyToken.Code).First(prevMoneyToken).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(moneyToken).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain money token's single column value in the database
func (repo *MoneyTokenRepository) UpdateValue(moneyToken *entity.MoneyToken, columnName string, columnValue interface{}) error {

	prevMoneyToken := new(entity.MoneyToken)
	err := repo.conn.Model(prevMoneyToken).Where("code = ?", moneyToken.Code).First(prevMoneyToken).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.MoneyToken{}).Where("code = ?", moneyToken.Code).Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain money token from the database using an identifier.
// In Delete() code is only used as a key
func (repo *MoneyTokenRepository) Delete(identifier string) (*entity.MoneyToken, error) {
	moneyToken := new(entity.MoneyToken)
	err := repo.conn.Model(moneyToken).Where("code = ?", identifier).First(moneyToken).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(moneyToken)
	return moneyToken, nil
}

// DeleteMultiple is a method that deletes multiple money tokens from the database using the identifier.
// In DeleteMultiple() sender_id is only used as a key
func (repo *MoneyTokenRepository) DeleteMultiple(identifier string) ([]*entity.MoneyToken, error) {
	var moneyTokens []*entity.MoneyToken
	err := repo.conn.Model(entity.MoneyToken{}).Where("sender_id = ?", identifier).Find(&moneyTokens).Error

	if err != nil {
		return nil, err
	}

	if len(moneyTokens) == 0 {
		return nil, errors.New("no money token for the provided identifier")
	}

	repo.conn.Model(entity.MoneyToken{}).Where("sender_id = ?", identifier).Delete(entity.MoneyToken{})
	return moneyTokens, nil
}

// IsUnique is a method that determines whether a certain column value is unique in the money tokens table
func (repo *MoneyTokenRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.MoneyToken{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
