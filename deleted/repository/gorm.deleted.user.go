package repository

import (
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
