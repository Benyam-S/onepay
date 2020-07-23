package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/entity"
	"github.com/google/uuid"
)

// AddAPIToken is a method that adds a new api token to the system using the api client
func (service *Service) AddAPIToken(apiToken *api.Token, apiClient *api.Client, opUser *entity.User) error {

	apiToken.AccessToken = "OP_Token-" + uuid.Must(uuid.NewRandom()).String()
	apiToken.APIKey = apiClient.APIKey
	apiToken.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()
	apiToken.DailyExpiration = time.Now().Unix()
	apiToken.UserID = opUser.UserID

	err := service.apiTokenRepo.Create(apiToken)
	if err != nil {
		return errors.New("unable to add new api token")
	}
	return nil
}

// FindAPIToken is a method that find and returns an api token for the given identifier
func (service *Service) FindAPIToken(identifier string) (*api.Token, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	apiToken, err := service.apiTokenRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("api token not found")
	}
	return apiToken, nil
}

// SearchAPIToken is a method that searchs and returns a set of tokens for the given identifier
func (service *Service) SearchAPIToken(identifier string) ([]*api.Token, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	apiTokens, err := service.apiTokenRepo.Search(identifier)
	if err != nil {
		return nil, errors.New("no api token found for the provided identifier")
	}

	return apiTokens, nil
}

// ValidateAPIToken is a method that checks if the api token is valid and have a valid api client
func (service *Service) ValidateAPIToken(apiToken *api.Token) error {

	if time.Now().Unix() > apiToken.ExpiresAt {
		return errors.New("invalid token, api token has expired")
	}

	// apiToken can be deactivated when a user logs out
	if apiToken.Deactivated {
		return errors.New("deactivated api token")
	}

	_, err := service.apiClientRepo.Find(apiToken.APIKey)
	if err != nil {
		return errors.New("api token not found")
	}

	return nil
}

// UpdateAPIToken is a method that updates a certain's api token
func (service *Service) UpdateAPIToken(apiToken *api.Token) error {

	err := service.apiTokenRepo.Update(apiToken)
	if err != nil {
		return errors.New("unable to update api token")
	}
	return nil
}

// DeleteAPIToken is a method that deletes an api token from the system
func (service *Service) DeleteAPIToken(identifier string) (*api.Token, error) {

	apiToken, err := service.apiTokenRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to deleted api token")
	}
	return apiToken, nil
}

// DeleteAPITokens is a method that deletes a set of api tokens from the system
func (service *Service) DeleteAPITokens(identifier string) ([]*api.Token, error) {

	apiTokens, err := service.apiTokenRepo.DeleteMultiple(identifier)
	if err != nil {
		return nil, errors.New("unable to deleted api tokens")
	}
	return apiTokens, nil
}
