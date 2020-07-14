package user

import (
	"net/http"

	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/entity"
)

// IService is an interface that defines all the service methods of a user struct
type IService interface {
	AddUser(opUser *entity.User, opPassword *entity.UserPassword) error
	FindUser(identifier string) (*entity.User, error)
	ValidateUserProfile(opUser *entity.User) entity.ErrMap
	UpdateUser(User *entity.User) error
	UpdateUserSingleValue(userID, columnName string, columnValue interface{}) error
	DeleteUser(userID string) (*entity.User, error)

	VerifyUserPassword(opPassword *entity.UserPassword, verifyPassword string) error
	FindPassword(identifier string) (*entity.UserPassword, error)
	UpdatePassword(opPassword *entity.UserPassword) error
	DeletePassword(identifier string) (*entity.UserPassword, error)

	AddSession(opClientSession *session.ClientSession, opUser *entity.User, r *http.Request) error
	FindSession(identifier string) ([]*session.ServerSession, error)
	UpdateSession(opServerSession *session.ServerSession) error
	DeleteSession(identifier string) (*session.ServerSession, error)

	AddAPIClient(apiClient *api.Client, opUser *entity.User) error
	FindAPIClient(identifier, clientType string) ([]*api.Client, error)
	UpdateAPIClient(apiClient *api.Client) error
	DeleteAPIClient(identifier string) (*api.Client, error)
	DeleteAPIClients(identifier string) ([]*api.Client, error)

	AddAPIToken(apiToken *api.Token, apiClient *api.Client, opUser *entity.User) error
	FindAPIToken(identifier string) ([]*api.Token, error)
	ValidateAPIToken(apiToken *api.Token) error
	UpdateAPIToken(apiToken *api.Token) error
	DeleteAPIToken(identifier string) (*api.Token, error)
	DeleteAPITokens(identifier string) ([]*api.Token, error)
}
