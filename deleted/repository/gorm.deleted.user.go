package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// DeletedUserRepository is a type that defines a repository for deleted user
type DeletedUserRepository struct {
	conn *gorm.DB
}

// NewDeletedUserRepository is a function that returns a new deleted user repository
func NewDeletedUserRepository(connection *gorm.DB) deleted.IDeletedUserRepository {
	return &DeletedUserRepository{conn: connection}
}

// Create is a method that adds a deleted user to the database
func (repo *DeletedUserRepository) Create(deletedOPUser *entity.DeletedUser) error {

	err := repo.conn.Create(deletedOPUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain deleted user from the database using an identifier,
// also Find() uses user_id as a key for selection
func (repo *DeletedUserRepository) Find(identifier string) (*entity.DeletedUser, error) {
	deletedOPUser := new(entity.DeletedUser)
	err := repo.conn.Model(deletedOPUser).
		Where("user_id = ? ", identifier).
		First(deletedOPUser).Error

	if err != nil {
		return nil, err
	}
	return deletedOPUser, nil
}

// Search is a method that search and returns a set of deleted users from the database using an identifier.
func (repo *DeletedUserRepository) Search(key string, pageNum int64, columns ...string) []*entity.DeletedUser {

	var deletedOPUsers []*entity.DeletedUser
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		// modifying the key so that it can match the database phone number values
		if column == "phone_number" {
			splitKey := strings.Split(key, "")
			if splitKey[0] == "0" {
				modifiedKey := "+251" + strings.Join(splitKey[1:], "")
				whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
				sqlValues = append(sqlValues, modifiedKey)
				continue
			}
		}
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_users WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedOPUsers)

	return deletedOPUsers
}

// SearchWRegx is a method that searchs and returns set of deleted users limited to the key identifier and page number using regular expersions
func (repo *DeletedUserRepository) SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedUser {
	var deletedOPUsers []*entity.DeletedUser
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_users WHERE "+strings.Join(whereStmt, "||")+" ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedOPUsers)

	return deletedOPUsers
}

// All is a method that returns all the deleted users from the database limited with the pageNum
func (repo *DeletedUserRepository) All(pageNum int64) []*entity.DeletedUser {

	var deletedOPUsers []*entity.DeletedUser
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_users ORDER BY first_name ASC LIMIT ?, 30", limit).Scan(&deletedOPUsers)
	return deletedOPUsers
}

// Update is a method that updates a certain deleted user value in the database
func (repo *DeletedUserRepository) Update(deletedOPUser *entity.DeletedUser) error {

	prevOPUser := new(entity.DeletedUser)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", deletedOPUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedOPUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain deleted user from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *DeletedUserRepository) Delete(identifier string) (*entity.DeletedUser, error) {
	deletedOPUser := new(entity.DeletedUser)
	err := repo.conn.Model(deletedOPUser).Where("user_id = ?", identifier).First(deletedOPUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedOPUser)
	return deletedOPUser, nil
}
