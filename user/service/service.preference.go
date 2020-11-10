package service

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Benyam-S/onepay/entity"
)

// FindUserPreference is a method that find and return a user preference that matchs the identifier value
func (service *Service) FindUserPreference(identifier string) (*entity.UserPreference, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	userPreference, err := service.preferenceRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("user preference not found")
	}

	return userPreference, nil
}

// ValidateUserPreference is a method that validates a user preference value according to the give column
func (service *Service) ValidateUserPreference(columnName, columValue string) (interface{}, error) {

	// Column name screeing
	validColumnNames := []string{"two_step_verification"}
	isValidColumnName := false

	for _, validColumnName := range validColumnNames {
		if columnName == validColumnName {
			isValidColumnName = true
			break
		}
	}

	if !isValidColumnName {
		return nil, errors.New("invalid column used")
	}

	if columnName == "two_step_verification" {
		value, err := strconv.ParseBool(columValue)
		if err != nil {
			return nil, errors.New("invalid value used")
		}

		return value, nil
	}

	return columValue, nil
}

// UpdateUserPreference is a method that updates a certain user's preference
func (service *Service) UpdateUserPreference(userPreference *entity.UserPreference) error {

	err := service.preferenceRepo.Update(userPreference)
	if err != nil {
		return errors.New("unable to update user preference")
	}
	return nil
}

// UpdateUserPreferenceSingleValue is a method that updates a single column entry of a user preference
func (service *Service) UpdateUserPreferenceSingleValue(userID, columnName string, columnValue interface{}) error {

	userPreference := entity.UserPreference{UserID: userID}
	err := service.preferenceRepo.UpdateValue(&userPreference, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update user preference")
	}

	return nil
}

// DeleteUserPreference is a method that deletes a certain user's preference
func (service *Service) DeleteUserPreference(identifier string) (*entity.UserPreference, error) {

	userPreference, err := service.preferenceRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete user preference")
	}
	return userPreference, nil
}
