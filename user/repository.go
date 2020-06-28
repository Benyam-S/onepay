package user

import (
	"github.com/Benyam-S/onepay/api"
	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/entity"
)

// IUserRepository is an interface that defines all the repository methods of a user struct
type IUserRepository interface {
	Create(newOPUser *entity.User) error
	Find(identifier string) (*entity.User, error)
	Update(opUser *entity.User) error
	Delete(identifier string) (*entity.User, error)
	CountUsers() int
	IsUnique(columnName string, columnValue interface{}) bool
}

// IPasswordRepository is an interface that defines all the repository methods of a user's password struct
type IPasswordRepository interface {
	Create(newOPPassword *entity.UserPassword) error
	Find(identifier string) (*entity.UserPassword, error)
	Update(opPassword *entity.UserPassword) error
	Delete(identifier string) (*entity.UserPassword, error)
}

// ISessionRepository is an interface that defines all the repository methods of a user's server side session struct
type ISessionRepository interface {
	Create(newOPSession *session.ServerSession) error
	Find(identifier string) ([]*session.ServerSession, error)
	Update(opSession *session.ServerSession) error
	Delete(identifier string) (*session.ServerSession, error)
}

// IAPIClientRepository is an interface that defines all the repository methods of an api client struct
type IAPIClientRepository interface {
	Create(newAPIClient *api.Client) error
	Find(identifier string) ([]*api.Client, error)
	Update(apiClient *api.Client) error
	Delete(identifier string) (*api.Client, error)
}

// IAPITokenRepository is an interface that defines all the repository methods of an api token struct
type IAPITokenRepository interface {
	Create(newAPIToken *api.Token) error
	Find(identifier string) ([]*api.Token, error)
	Update(apiToken *api.Token) error
	Delete(identifier string) (*api.Token, error)
}
