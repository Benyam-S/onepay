package repository

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// PasswordRepository is a type that defines a user's password repository
type PasswordRepository struct {
	conn *gorm.DB
}

// NewPasswordRepository is a function that returns a new user's password repository
func NewPasswordRepository(connection *gorm.DB) user.IPasswordRepository {
	return &PasswordRepository{conn: connection}
}

// Create is a method that adds a new user password to the database
func (repo *PasswordRepository) Create(newOPPassword *entity.UserPassword) error {

	err := repo.conn.Create(newOPPassword).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain user's password from the database using an identifier.
// In Find() user_id is only used as a key
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
func (repo *PasswordRepository) Update(opPassword *entity.UserPassword) error {

	prevOPPassword := new(entity.UserPassword)
	err := repo.conn.Model(prevOPPassword).Where("user_id = ?", opPassword.UserID).First(prevOPPassword).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(opPassword).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user's password from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *PasswordRepository) Delete(identifier string) (*entity.UserPassword, error) {
	opPassword := new(entity.UserPassword)
	err := repo.conn.Model(opPassword).Where("user_id = ?", identifier).First(opPassword).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opPassword)
	return opPassword, nil
}
