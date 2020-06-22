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

// PasswordRepository is a type that defines a user's password repository
type PasswordRepository struct {
	conn *gorm.DB
}

// NewUserRepository is a function that returns a new user repository
func NewUserRepository(connection *gorm.DB) user.IUserRepository {
	return &UserRepository{conn: connection}
}

// NewPasswordRepository is a function that returns a new user's password repository
func NewPasswordRepository(connection *gorm.DB) user.IPasswordRepository {
	return &PasswordRepository{conn: connection}
}

/* ================================= UserRepository Methods ================================= */

// Create is a method that adds a new user to the database
func (repo *UserRepository) Create(newOPUser *entity.User) (*entity.User, error) {
	totalNumOfUsers := repo.CountUser()
	newOPUser.UserID = fmt.Sprintf("OP%s%d", tools.RandomStringGN(7), totalNumOfUsers+1)

	for !repo.IsUnique("user_id", newOPUser.UserID) {
		totalNumOfUsers++
		newOPUser.UserID = fmt.Sprintf("OP%s%d", tools.RandomStringGN(7), totalNumOfUsers+1)
	}

	err := repo.conn.Create(newOPUser).Error
	if err != nil {
		return nil, err
	}
	return newOPUser, nil
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
func (repo *UserRepository) Update(opUser *entity.User) (*entity.User, error) {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ?", opUser.UserID).Find(prevOPUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Save(opUser)
	return opUser, nil
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

/* ================================= PasswordRepository Methods ================================= */

// Create is a method that adds a new user password to the database
func (repo *PasswordRepository) Create(newOPPassword *entity.UserPassword) (*entity.UserPassword, error) {

	err := repo.conn.Create(newOPPassword).Error
	if err != nil {
		return nil, err
	}
	return newOPPassword, nil
}

// Find is a method that finds a certain user's password from the database using an identifier.
// In Find() user_id is only used as an key
func (repo *PasswordRepository) Find(identifier string) (*entity.UserPassword, error) {
	opPassword := new(entity.UserPassword)
	err := repo.conn.Model(opPassword).
		Where("user_id = ?", identifier).
		First(opPassword).Error

	if err != nil {
		return nil, err
	}
	return opPassword, nil
}

// Update is a method that updates a certain user's password value in the database
func (repo *PasswordRepository) Update(opPassword *entity.UserPassword) (*entity.UserPassword, error) {

	prevOPPassword := new(entity.UserPassword)
	err := repo.conn.Model(opPassword).Where("user_id = ?", opPassword.UserID).Find(prevOPPassword).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Save(opPassword)
	return opPassword, nil
}

// Delete is a method that deletes a certain user's password from the database using an identifier.
// In Delete() user_id is only used as an key
func (repo *PasswordRepository) Delete(identifier string) (*entity.UserPassword, error) {
	opPassword := new(entity.UserPassword)
	err := repo.conn.Model(opPassword).Where("user_id = ?", identifier).First(opPassword).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opPassword)
	return opPassword, nil
}
