package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/Benyam-S/onepay/tools"

	"github.com/Benyam-S/onepay/accountprovider"
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// AccountProviderRepository is a type that defines a account provider repository
type AccountProviderRepository struct {
	conn *gorm.DB
}

// NewAccountProviderRepository is a function that returns a new account provider repository
func NewAccountProviderRepository(connection *gorm.DB) accountprovider.IAccountProviderRepository {
	return &AccountProviderRepository{conn: connection}
}

// Create is a method that adds a new account provider to the database
func (repo *AccountProviderRepository) Create(newAccountProvider *entity.AccountProvider) error {
	newAccountProvider.ID = "OP_AP-" + strings.ToUpper(tools.GenerateRandomString(9))

	for !repo.IsUnique("id", newAccountProvider.ID) {
		newAccountProvider.ID = "OP_AP-" + strings.ToUpper(tools.GenerateRandomString(9))
	}

	err := repo.conn.Create(newAccountProvider).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain account provider from the database using an identifier.
// In Find() id is only used as a key
func (repo *AccountProviderRepository) Find(identifier string) (*entity.AccountProvider, error) {
	accountProvider := new(entity.AccountProvider)
	err := repo.conn.Model(accountProvider).
		Where("id = ?", identifier).
		First(accountProvider).Error

	if err != nil {
		return nil, err
	}
	return accountProvider, nil
}

// All is a method that returns all the account providers found in the database
func (repo *AccountProviderRepository) All() []*entity.AccountProvider {
	var accountProviders []*entity.AccountProvider
	err := repo.conn.Model(entity.AccountProvider{}).Find(&accountProviders).Error

	if err != nil {
		return []*entity.AccountProvider{}
	}
	return accountProviders
}

// Search is a method that searchs for an account provider that match the column name and value.
func (repo *AccountProviderRepository) Search(columnName string, columnValue interface{}) []*entity.AccountProvider {
	var accountProviders []*entity.AccountProvider
	err := repo.conn.Model(entity.AccountProvider{}).
		Where(columnName+" = ?", columnValue).
		Find(&accountProviders).Error

	if err != nil {
		return []*entity.AccountProvider{}
	}
	return accountProviders
}

// SearchMultiple is a method that search and returns a set of account providers that is limited to the provided identifier
func (repo *AccountProviderRepository) SearchMultiple(key, orderBy string, pageNum int64, columns ...string) ([]*entity.AccountProvider, int64) {

	var accountProviders []*entity.AccountProvider
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	repo.conn.Raw("SELECT COUNT(*) FROM account_providers WHERE ("+strings.Join(whereStmt, "||")+")", sqlValues...).Count(&count)
	repo.conn.Raw("SELECT * FROM account_providers WHERE ("+strings.Join(whereStmt, "||")+")", sqlValues...).Order(orderBy + " DESC").Limit(5).Offset(pageNum * 5).Scan(&accountProviders)

	var pageCount int64 = int64(math.Ceil(count / 5.0))
	return accountProviders, pageCount
}

// Update is a method that updates a certain account provider value in the database
func (repo *AccountProviderRepository) Update(accountProvider *entity.AccountProvider) error {

	prevAccountProvider := new(entity.AccountProvider)
	err := repo.conn.Model(prevAccountProvider).Where("id = ?", accountProvider.ID).
		First(prevAccountProvider).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	accountProvider.CreatedAt = prevAccountProvider.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(accountProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain account provider's single column value in the database
func (repo *AccountProviderRepository) UpdateValue(accountProvider *entity.AccountProvider, columnName string, columnValue interface{}) error {

	prevAccountProvider := new(entity.AccountProvider)
	err := repo.conn.Model(prevAccountProvider).Where("id = ?", accountProvider.ID).
		First(prevAccountProvider).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.AccountProvider{}).Where("id = ?", accountProvider.ID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain account provider from the database using an identifier.
// In Delete() id is only used as a key
func (repo *AccountProviderRepository) Delete(identifier string) (*entity.AccountProvider, error) {
	accountProvider := new(entity.AccountProvider)
	err := repo.conn.Model(accountProvider).Where("id = ?", identifier).First(accountProvider).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(accountProvider)
	return accountProvider, nil
}

// IsUnique is a method that determines whether a certain column value is unique in the account providers table
func (repo *AccountProviderRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.AccountProvider{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
