package service

import (
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
	return err
}
