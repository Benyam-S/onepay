package service

import (
	"errors"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// AddAPIClient is a method that adds a new api client to the system using the one pay user
func (service *Service) AddAPIClient(apiClient *api.Client, opUser *entity.User) error {

	apiClient.ClientUserID = opUser.UserID
	apiClient.APISecret = tools.GenerateRandomString(20)

	err := service.apiClientRepo.Create(apiClient)
	return err
}

// FindAPIClient is a method that finds a client from the system using the given identifier and client type
func (service *Service) FindAPIClient(identifier, clientType string) ([]*api.Client, error) {

	apiClientsUnFiltered, err := service.apiClientRepo.Find(identifier)
	if err != nil {
		return nil, err
	}

	apiClientsFiltered := make([]*api.Client, 0)
	for _, client := range apiClientsUnFiltered {
		if client.Type == clientType {
			apiClientsFiltered = append(apiClientsFiltered, client)
		}
	}

	if len(apiClientsFiltered) == 0 {
		return nil, errors.New("no api client found for the provided identifier and filter")
	}

	return apiClientsFiltered, nil
}

// UpdateAPIClient is a method that updates a certain's api client
func (service *Service) UpdateAPIClient(apiClient *api.Client) error {
	return service.apiClientRepo.Update(apiClient)
}

// DeleteAPIClient is a method that deletes a certain's api client using the identifier
func (service *Service) DeleteAPIClient(identifier string) (*api.Client, error) {
	return service.apiClientRepo.Delete(identifier)
}

// DeleteAPIClients is a method that deletes a set of api client from the system that matchs the given identifier
func (service *Service) DeleteAPIClients(identifier string) ([]*api.Client, error) {
	return service.apiClientRepo.DeleteMultiple(identifier)
}
