package repository

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// DeletedLinkedAccountRepository is a type that defines a repository for deleted linked accounts
type DeletedLinkedAccountRepository struct {
	conn *gorm.DB
}

// NewDeletedLinkedAccountRepository is a function that returns a new deleted linked account repository
func NewDeletedLinkedAccountRepository(connection *gorm.DB) deleted.IDeletedLinkedAccountRepository {
	return &DeletedLinkedAccountRepository{conn: connection}
}

// Create is a method that adds a deleted linked account to the database
func (repo *DeletedLinkedAccountRepository) Create(deletedLinkedAccount *entity.DeletedLinkedAccount) error {

	err := repo.conn.Create(deletedLinkedAccount).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain deleted linked account from the database using an identifier.
// In Find() id is only used as a key
func (repo *DeletedLinkedAccountRepository) Find(identifier string) (*entity.DeletedLinkedAccount, error) {
	deletedLinkedAccount := new(entity.DeletedLinkedAccount)
	err := repo.conn.Model(deletedLinkedAccount).
		Where("id = ?", identifier).
		First(deletedLinkedAccount).Error

	if err != nil {
		return nil, err
	}
	return deletedLinkedAccount, nil
}

// Search is a method that searchs for a linked account that match the column name and value.
func (repo *DeletedLinkedAccountRepository) Search(colunmName string, columnValue interface{}) []*entity.DeletedLinkedAccount {

	var deletedLinkedAccounts []*entity.DeletedLinkedAccount
	err := repo.conn.Model(entity.DeletedLinkedAccount{}).
		Where(colunmName+" = ?", columnValue).
		Find(&deletedLinkedAccounts).Error

	if err != nil {
		return []*entity.DeletedLinkedAccount{}
	}
	return deletedLinkedAccounts
}

// SearchMultiple is a method that search and returns a set of deleted linked accounts from the database using an identifier.
func (repo *DeletedLinkedAccountRepository) SearchMultiple(key string, pageNum int64, columns ...string) ([]*entity.DeletedLinkedAccount, int64) {

	var deletedLinkedAccounts []*entity.DeletedLinkedAccount
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	repo.conn.Raw("SELECT COUNT(*) FROM deleted_linked_accounts WHERE ("+strings.Join(whereStmt, "||")+")", sqlValues...).
		Count(&count)
	repo.conn.Raw("SELECT * FROM deleted_linked_accounts WHERE ("+strings.Join(whereStmt, "||")+")", sqlValues...).
		Order("id ASC").Limit(30).Offset(pageNum * 30).Scan(&deletedLinkedAccounts)

	var pageCount int64 = int64(math.Ceil(count / 30.0))
	return deletedLinkedAccounts, pageCount
}

// Update is a method that updates a certain deleted linked account value in the database
func (repo *DeletedLinkedAccountRepository) Update(deletedLinkedAccount *entity.DeletedLinkedAccount) error {

	prevLinkedAccount := new(entity.DeletedLinkedAccount)
	err := repo.conn.Model(prevLinkedAccount).Where("id = ?", deletedLinkedAccount.ID).First(prevLinkedAccount).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedLinkedAccount).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain deleted linked account from the database using an identifier.
// In Delete() id is only used as a key
func (repo *DeletedLinkedAccountRepository) Delete(identifier string) (*entity.DeletedLinkedAccount, error) {
	deletedLinkedAccount := new(entity.DeletedLinkedAccount)
	err := repo.conn.Model(deletedLinkedAccount).Where("id = ?", identifier).
		First(deletedLinkedAccount).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedLinkedAccount)
	return deletedLinkedAccount, nil
}

// DeleteMultiple is a method that deletes multiple deleted linked accounts from the database using the identifier.
// In DeleteMultiple() user_id is only used as a key
func (repo *DeletedLinkedAccountRepository) DeleteMultiple(identifier string) ([]*entity.DeletedLinkedAccount, error) {
	var deletedLinkedAccounts []*entity.DeletedLinkedAccount
	err := repo.conn.Model(entity.DeletedLinkedAccount{}).Where("user_id = ?", identifier).Find(&deletedLinkedAccounts).Error

	if err != nil {
		return nil, err
	}

	if len(deletedLinkedAccounts) == 0 {
		return nil, errors.New("no deleted linked account for the provided identifier")
	}

	repo.conn.Model(entity.DeletedLinkedAccount{}).Where("user_id = ?", identifier).Delete(entity.DeletedLinkedAccount{})
	return deletedLinkedAccounts, nil
}

// IsUnique is a method that determines whether a certain column value is unique in the deleted linked accounts table
// Is is used to create a unique ID for the deleted linked account since duplication may occur
func (repo *DeletedLinkedAccountRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.DeletedLinkedAccount{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
