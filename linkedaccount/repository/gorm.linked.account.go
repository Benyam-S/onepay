package repository

import (
	"errors"
	"fmt"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/linkedaccount"
	"github.com/Benyam-S/onepay/tools"
	"github.com/jinzhu/gorm"
)

// LinkedAccountRepository is a type that defines a linked account repository
type LinkedAccountRepository struct {
	conn *gorm.DB
}

// NewLinkedAccountRepository is a function that returns a new linked account repository
func NewLinkedAccountRepository(connection *gorm.DB) linkedaccount.ILinkedAccountRepository {
	return &LinkedAccountRepository{conn: connection}
}

// Create is a method that adds a new linked account to the database
func (repo *LinkedAccountRepository) Create(newLinkedAccount *entity.LinkedAccount) error {

	newLinkedAccount.ID = fmt.Sprintf("OP_LA-%s%s", newLinkedAccount.UserID[3:]+"_", tools.GenerateRandomString(5))

	for !repo.IsUnique("id", newLinkedAccount.ID) {
		newLinkedAccount.ID = fmt.Sprintf("OP_LA-%s%s", newLinkedAccount.UserID[3:]+"_", tools.GenerateRandomString(5))
	}

	err := repo.conn.Create(newLinkedAccount).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain linked account from the database using an identifier.
// In Find() id is only used as a key
func (repo *LinkedAccountRepository) Find(identifier string) (*entity.LinkedAccount, error) {
	linkedAccount := new(entity.LinkedAccount)
	err := repo.conn.Model(linkedAccount).
		Where("id = ?", identifier).
		First(linkedAccount).Error

	if err != nil {
		return nil, err
	}
	return linkedAccount, nil
}

// Search is a method that searchs for a linked account that match the column name and value.
func (repo *LinkedAccountRepository) Search(colunmName string, columnValue interface{}) []*entity.LinkedAccount {
	var linkedAccounts []*entity.LinkedAccount
	err := repo.conn.Model(entity.LinkedAccount{}).
		Where(colunmName+" = ?", columnValue).
		Find(&linkedAccounts).Error

	if err != nil {
		return []*entity.LinkedAccount{}
	}
	return linkedAccounts
}

// Update is a method that updates a certain linked account value in the database
func (repo *LinkedAccountRepository) Update(linkedAccount *entity.LinkedAccount) error {

	prevLinkedAccount := new(entity.LinkedAccount)
	err := repo.conn.Model(prevLinkedAccount).Where("id = ?", linkedAccount.ID).First(prevLinkedAccount).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(linkedAccount).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain linked account's single column value in the database
func (repo *LinkedAccountRepository) UpdateValue(linkedAccount *entity.LinkedAccount, columnName string, columnValue interface{}) error {

	prevLinkedAccount := new(entity.LinkedAccount)
	err := repo.conn.Model(prevLinkedAccount).Where("id = ?", linkedAccount.ID).First(prevLinkedAccount).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.LinkedAccount{}).Where("id = ?", linkedAccount.ID).Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain linked account from the database using an identifier.
// In Delete() id is only used as a key
func (repo *LinkedAccountRepository) Delete(identifier string) (*entity.LinkedAccount, error) {
	linkedAccount := new(entity.LinkedAccount)
	err := repo.conn.Model(linkedAccount).Where("id = ?", identifier).First(linkedAccount).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(linkedAccount)
	return linkedAccount, nil
}

// DeleteMultiple is a method that deletes multiple linked accounts from the database using the identifier.
// In DeleteMultiple() user_id is only used as a key
func (repo *LinkedAccountRepository) DeleteMultiple(identifier string) ([]*entity.LinkedAccount, error) {
	var linkedAccounts []*entity.LinkedAccount
	err := repo.conn.Model(entity.LinkedAccount{}).Where("user_id = ?", identifier).Find(&linkedAccounts).Error

	if err != nil {
		return nil, err
	}

	if len(linkedAccounts) == 0 {
		return nil, errors.New("no linked account for the provided identifier")
	}

	repo.conn.Model(entity.LinkedAccount{}).Where("user_id = ?", identifier).Delete(entity.LinkedAccount{})
	return linkedAccounts, nil
}

// IsUnique is a method that determines whether a certain column value is unique in the linked accounts table
func (repo *LinkedAccountRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.LinkedAccount{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
