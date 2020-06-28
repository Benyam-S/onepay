package repository

import (
	"errors"
	"fmt"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// APIClientRepository is a type that defines an api client repository
type APIClientRepository struct {
	conn *gorm.DB
}

// NewAPIClientRepository is a function that returns a new api client repository
func NewAPIClientRepository(connection *gorm.DB) user.IAPIClientRepository {
	return &APIClientRepository{conn: connection}
}

// Create is a method that adds a new api client to the database
func (repo *APIClientRepository) Create(newAPIClient *api.Client) error {
	totalNumOfClients := repo.CountClients()
	newAPIClient.APIKey = fmt.Sprintf("OP_API-%s%d", tools.RandomStringGN(10), totalNumOfClients+1)

	for !repo.IsUnique("api_key", newAPIClient.APIKey) {
		totalNumOfClients++
		newAPIClient.APIKey = fmt.Sprintf("OP_API-%s%d", tools.RandomStringGN(10), totalNumOfClients+1)
	}

	err := repo.conn.Create(newAPIClient).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that find an api client from the database using an identifier.
// In Find() client_user_id and api_key can be used as an key
func (repo *APIClientRepository) Find(identifier string) ([]*api.Client, error) {
	var apiClients []*api.Client
	err := repo.conn.Model(api.Client{}).
		Where("client_user_id = ? || api_key = ?", identifier, identifier).
		Find(&apiClients).Error

	if err != nil {
		return nil, err
	}

	if len(apiClients) == 0 {
		return nil, errors.New("no available api client for the provided identifier")
	}
	return apiClients, nil
}

// Update is a method that updates an api client value in the database
func (repo *APIClientRepository) Update(apiClient *api.Client) error {

	prevAPIClient := new(api.Client)
	err := repo.conn.Model(prevAPIClient).Where("api_key = ?", apiClient.APIKey).First(prevAPIClient).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(api.Client{}).Where("api_key = ?", apiClient.APIKey).Update(apiClient).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain api client from the database using an identifier.
// In Delete() api_key is only used as an key
func (repo *APIClientRepository) Delete(identifier string) (*api.Client, error) {
	apiClient := new(api.Client)
	err := repo.conn.Model(apiClient).Where("api_key = ?", identifier).First(apiClient).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(apiClient)
	return apiClient, nil
}

// CountClients is a method that counts the api clients in the database
func (repo *APIClientRepository) CountClients() int {
	var totalNumOfClients int
	repo.conn.Model(&api.Client{}).Unscoped().Count(&totalNumOfClients)
	return totalNumOfClients
}

// IsUnique is a method that determines whether a certain column value is unique in the api clients table
func (repo *APIClientRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&api.Client{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}
