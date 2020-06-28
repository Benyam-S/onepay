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
	apiClient.APISecret = tools.RandomStringGN(20)

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
