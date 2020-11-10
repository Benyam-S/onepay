package repository

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// PreferenceRepository is a type that defines a user preference repository
type PreferenceRepository struct {
	conn *gorm.DB
}

// NewPreferenceRepository is a function that returns a new user preference repository
func NewPreferenceRepository(connection *gorm.DB) user.IPreferenceRepository {
	return &PreferenceRepository{conn: connection}
}

// Create is a method that adds a new user preference to the database
func (repo *PreferenceRepository) Create(newUserPreference *entity.UserPreference) error {
	err := repo.conn.Create(newUserPreference).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain user preference from the database using an identifier,
// also Find() uses only user_id as a key for selection
func (repo *PreferenceRepository) Find(identifier string) (*entity.UserPreference, error) {
	userPreference := new(entity.UserPreference)
	err := repo.conn.Model(userPreference).
		Where("user_id = ?", identifier).
		First(userPreference).Error

	if err != nil {
		return nil, err
	}
	return userPreference, nil
}

// Update is a method that updates a certain user preference value in the database
func (repo *PreferenceRepository) Update(userPreference *entity.UserPreference) error {

	preUserPreference := new(entity.UserPreference)
	err := repo.conn.Model(preUserPreference).Where("user_id = ?", userPreference.UserID).First(preUserPreference).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(userPreference).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain user preference single column value in the database
func (repo *PreferenceRepository) UpdateValue(userPreference *entity.UserPreference, columnName string, columnValue interface{}) error {

	preUserPreference := new(entity.UserPreference)
	err := repo.conn.Model(preUserPreference).Where("user_id = ?", userPreference.UserID).First(preUserPreference).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.UserPreference{}).Where("user_id = ?", userPreference.UserID).Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user preference from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *PreferenceRepository) Delete(identifier string) (*entity.UserPreference, error) {
	userPreference := new(entity.UserPreference)
	err := repo.conn.Model(userPreference).Where("user_id = ?", identifier).First(userPreference).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(userPreference)
	return userPreference, nil
}
