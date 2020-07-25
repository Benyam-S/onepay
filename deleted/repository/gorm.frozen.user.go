package repository

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// FrozenUserRepository is a type that defines a repository for frozen user
type FrozenUserRepository struct {
	conn *gorm.DB
}

// NewFrozenUserRepository is a function that returns a new frozen user repository
func NewFrozenUserRepository(connection *gorm.DB) deleted.IFrozenUserRepository {
	return &FrozenUserRepository{conn: connection}
}

// Create is a method that adds a frozen user to the database
func (repo *FrozenUserRepository) Create(frozenOPUser *entity.FrozenUser) error {

	err := repo.conn.Create(frozenOPUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain frozen user from the database using an identifier,
// also Find() uses user_id as a key for selection
func (repo *FrozenUserRepository) Find(identifier string) (*entity.FrozenUser, error) {
	frozenOPUser := new(entity.FrozenUser)
	err := repo.conn.Model(frozenOPUser).
		Where("user_id = ? ", identifier).
		First(frozenOPUser).Error

	if err != nil {
		return nil, err
	}
	return frozenOPUser, nil
}

// Search is a method that search and returns a set of frozen users from the database using an identifier.
func (repo *FrozenUserRepository) Search(key string, pageNum int64, columns ...string) []*entity.FrozenUser {

	var frozenOPUsers []*entity.FrozenUser
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM frozen_users WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY user_id ASC LIMIT ?, 30", sqlValues...).Scan(&frozenOPUsers)

	return frozenOPUsers
}

// All is a method that returns all the frozen users from the database limited with the pageNum
func (repo *FrozenUserRepository) All(pageNum int64) []*entity.FrozenUser {

	var frozenOPUsers []*entity.FrozenUser
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM frozen_users ORDER BY user_id ASC LIMIT ?, 30", limit).Scan(&frozenOPUsers)
	return frozenOPUsers
}

// Update is a method that updates a certain frozen user value in the database
func (repo *FrozenUserRepository) Update(frozenOPUser *entity.FrozenUser) error {

	prevOPUser := new(entity.FrozenUser)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", frozenOPUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	frozenOPUser.CreatedAt = prevOPUser.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(frozenOPUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain frozen user from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *FrozenUserRepository) Delete(identifier string) (*entity.FrozenUser, error) {
	frozenOPUser := new(entity.FrozenUser)
	err := repo.conn.Model(frozenOPUser).Where("user_id = ?", identifier).First(frozenOPUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(frozenOPUser)
	return frozenOPUser, nil
}
