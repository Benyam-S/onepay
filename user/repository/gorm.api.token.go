package repository

import (
	"errors"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// APITokenRepository is a type that defines an api token repository
type APITokenRepository struct {
	conn *gorm.DB
}

// NewAPITokenRepository is a function that returns a new api token repository
func NewAPITokenRepository(connection *gorm.DB) user.IAPITokenRepository {
	return &APITokenRepository{conn: connection}
}

// Create is a method that adds a new api token to the database
func (repo *APITokenRepository) Create(newAPIToken *api.Token) error {

	err := repo.conn.Create(newAPIToken).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that find an api token from the database using an identifier.
// In Find() access_token and api_key can be used as a key
func (repo *APITokenRepository) Find(identifier string) ([]*api.Token, error) {
	var apiTokens []*api.Token
	err := repo.conn.Model(api.Token{}).
		Where("access_token = ? || api_key = ?", identifier, identifier).
		Find(&apiTokens).Error

	if err != nil {
		return nil, err
	}

	if len(apiTokens) == 0 {
		return nil, errors.New("no available api token for the provided identifier")
	}
	return apiTokens, nil
}

// Update is a method that updates an api token value in the database
func (repo *APITokenRepository) Update(apiToken *api.Token) error {

	prevAPIToken := new(api.Token)
	err := repo.conn.Model(prevAPIToken).Where("access_token = ?", apiToken.AccessToken).First(prevAPIToken).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(api.Token{}).Where("access_token = ?", apiToken.AccessToken).Update(apiToken).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain api token from the database using an identifier.
// In Delete() access_token is only used as a key
func (repo *APITokenRepository) Delete(identifier string) (*api.Token, error) {
	apiToken := new(api.Token)
	err := repo.conn.Model(apiToken).Where("access_token = ?", identifier).First(apiToken).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(apiToken)
	return apiToken, nil
}

// DeleteMultiple is a method that deletes multiple api tokens from the database using the identifier.
// In DeleteMultiple() user_id is only as a key, we didn't used api_key as a key identifier because it may cause loss of user session data
func (repo *APITokenRepository) DeleteMultiple(identifier string) ([]*api.Token, error) {
	var apiTokens []*api.Token
	err := repo.conn.Model(api.Token{}).Where("user_id = ?", identifier).
		Find(&apiTokens).Error

	if err != nil {
		return nil, err
	}

	if len(apiTokens) == 0 {
		return nil, errors.New("no api token for the provided identifier")
	}

	repo.conn.Model(api.Token{}).Where("user_id = ?", identifier).Delete(api.Token{})
	return apiTokens, nil
}
