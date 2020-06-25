package repository

import (
	"fmt"

	"github.com/Benyam-S/onepay/tools"

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
	totalNumOfUsers := repo.CountUser()
	newOPUser.UserID = fmt.Sprintf("OP%s%d", tools.RandomStringGN(7), totalNumOfUsers+1)

	for !repo.IsUnique("user_id", newOPUser.UserID) {
		totalNumOfUsers++
		newOPUser.UserID = fmt.Sprintf("OP%s%d", tools.RandomStringGN(7), totalNumOfUsers+1)
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
	opUser := new(entity.User)
	err := repo.conn.Model(opUser).
		Where("user_id = ? || email = ? || phone_number = ?", identifier, identifier, identifier).
		First(opUser).Error

	if err != nil {
		return nil, err
	}
	return opUser, nil
}

// Update is a method that updates a certain user value in the database
func (repo *UserRepository) Update(opUser *entity.User) error {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ?", opUser.UserID).Find(prevOPUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(opUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user from the database using an identifier.
// In Delete() user_id is only used as an key
func (repo *UserRepository) Delete(identifier string) (*entity.User, error) {
	opUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ?", identifier).First(opUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opUser)
	return opUser, nil
}

// CountUser is a method that counts the user in the database
func (repo *UserRepository) CountUser() int {
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
