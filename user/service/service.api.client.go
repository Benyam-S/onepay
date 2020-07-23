package service

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// AddAPIClient is a method that adds a new api client to the system using the one pay user
func (service *Service) AddAPIClient(apiClient *api.Client, opUser *entity.User) error {

	apiClient.ClientUserID = opUser.UserID
	apiClient.APISecret = tools.GenerateRandomString(20)

	err := service.apiClientRepo.Create(apiClient)
	if err != nil {
		return errors.New("unable to add new api client")
	}
	return nil
}

// FindAPIClient is a method that finds a client from the system using the given identifier and client type
func (service *Service) FindAPIClient(identifier string) (*api.Client, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	apiClient, err := service.apiClientRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("api client not found")
	}

	return apiClient, nil
}

// SearchAPIClient is a method that searchs for clients from the system using the given identifier and client type
func (service *Service) SearchAPIClient(identifier, clientType string) ([]*api.Client, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	apiClientsUnFiltered, err := service.apiClientRepo.Search(identifier)
	if err != nil {
		return nil, errors.New("no api client found for the provided identifier and filter")
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

	err := service.apiClientRepo.Update(apiClient)
	if err != nil {
		return errors.New("unable to update api client")
	}
	return nil
}

// DeleteAPIClient is a method that deletes a certain's api client using the identifier
func (service *Service) DeleteAPIClient(identifier string) (*api.Client, error) {

	apiClient, err := service.apiClientRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete api client")
	}
	return apiClient, nil
}

// DeleteAPIClients is a method that deletes a set of api client from the system that matchs the given identifier
func (service *Service) DeleteAPIClients(identifier string) ([]*api.Client, error) {

	apiClients, err := service.apiClientRepo.DeleteMultiple(identifier)
	if err != nil {
		return nil, errors.New("unable to delete api clients")
	}
	return apiClients, nil
}
