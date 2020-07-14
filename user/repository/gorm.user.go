package repository

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// UserRepository is a type that defines a user repository
type UserRepository struct {
	conn *gorm.DB
}

// NewUserRepository is a function that returns a new user repository
func NewUserRepository(connection *gorm.DB) user.IUserRepository {
	return &UserRepository{conn: connection}
}

// Create is a method that adds a new user to the database
func (repo *UserRepository) Create(newOPUser *entity.User) error {
	totalNumOfUsers := repo.CountUsers()
	baseID := 1000101101010
	newOPUser.UserID = fmt.Sprintf("OP-%d", baseID+totalNumOfUsers)

	for !repo.IsUnique("user_id", newOPUser.UserID) {
		newOPUser.UserID = fmt.Sprintf("OP-%d", baseID+totalNumOfUsers)
		totalNumOfUsers++
	}

	err := repo.conn.Create(newOPUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain user from the database using an identifier,
// also Find() uses user_id, email, phone_number as a key for selection
func (repo *UserRepository) Find(identifier string) (*entity.User, error) {

	var modifiedIdentifier string
	splitedIdentifier := strings.Split(identifier, "")
	if splitedIdentifier[0] == "0" {
		modifiedIdentifier = "+251" + strings.Join(splitedIdentifier[1:], "")
	}

	opUser := new(entity.User)
	err := repo.conn.Model(opUser).
		Where("user_id = ? || email = ? || phone_number = ?", identifier, identifier, modifiedIdentifier).
		First(opUser).Error

	if err != nil {
		return nil, err
	}
	return opUser, nil
}

// Update is a method that updates a certain user value in the database
func (repo *UserRepository) Update(opUser *entity.User) error {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", opUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	if opUser.ProfilePic == "" {
		opUser.ProfilePic = prevOPUser.ProfilePic
	}
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(opUser).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain user's single column value in the database
func (repo *UserRepository) UpdateValue(opUser *entity.User, columnName string, columnValue interface{}) error {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", opUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.User{}).Where("user_id = ?", opUser.UserID).Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *UserRepository) Delete(identifier string) (*entity.User, error) {
	opUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ?", identifier).First(opUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opUser)
	return opUser, nil
}

// CountUsers is a method that counts the users in the database
func (repo *UserRepository) CountUsers() int {
	var totalNumOfUsers int
	repo.conn.Model(&entity.User{}).Count(&totalNumOfUsers)
	return totalNumOfUsers
}

// IsUnique is a method that determines whether a certain column value is unique in the user table
func (repo *UserRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.User{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
