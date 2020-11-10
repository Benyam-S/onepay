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
	FindAlsoWPhone(identifier, phoneNumber string) (*entity.User, error)
	Search(key string, pageNum int64, columns ...string) []*entity.User
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.User
	All(pageNum int64) []*entity.User
	Update(opUser *entity.User) error
	UpdateValue(opUser *entity.User, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.User, error)
	IsUnique(columnName string, columnValue interface{}) bool
	IsUniqueRexp(columnName string, columnPattern string) bool
}

// IPasswordRepository is an interface that defines all the repository methods of a user's password struct
type IPasswordRepository interface {
	Create(newOPPassword *entity.UserPassword) error
	Find(identifier string) (*entity.UserPassword, error)
	Update(opPassword *entity.UserPassword) error
	Delete(identifier string) (*entity.UserPassword, error)
}

// IPreferenceRepository is an interface that defines all the repository methods of a user preference struct
type IPreferenceRepository interface {
	Create(newUserPreference *entity.UserPreference) error
	Find(identifier string) (*entity.UserPreference, error)
	Update(userPreference *entity.UserPreference) error
	UpdateValue(userPreference *entity.UserPreference, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.UserPreference, error)
}

// ISessionRepository is an interface that defines all the repository methods of a user's server side session struct
type ISessionRepository interface {
	Create(newOPSession *session.ServerSession) error
	Find(identifier string) (*session.ServerSession, error)
	Search(identifier string) ([]*session.ServerSession, error)
	SearchMultiple(key string, pageNum int64, columns ...string) []*session.ServerSession
	SearchMultipleWRegx(key string, pageNum int64, columns ...string) []*session.ServerSession
	Update(opSession *session.ServerSession) error
	Delete(identifier string) (*session.ServerSession, error)
	DeleteMultiple(identifier string) ([]*session.ServerSession, error)
}

// IAPIClientRepository is an interface that defines all the repository methods of an api client struct
type IAPIClientRepository interface {
	Create(newAPIClient *api.Client) error
	Find(identifier string) (*api.Client, error)
	Search(identifier string) ([]*api.Client, error)
	SearchMultiple(key string, pageNum int64, columns ...string) []*api.Client
	All(pageNum int64) []*api.Client
	Update(apiClient *api.Client) error
	Delete(identifier string) (*api.Client, error)
	DeleteMultiple(identifier string) ([]*api.Client, error)
}

// IAPITokenRepository is an interface that defines all the repository methods of an api token struct
type IAPITokenRepository interface {
	Create(newAPIToken *api.Token) error
	Find(identifier string) (*api.Token, error)
	Search(identifier string) ([]*api.Token, error)
	SearchMultiple(key string, pageNum int64, columns ...string) []*api.Token
	Update(apiToken *api.Token) error
	Delete(identifier string) (*api.Token, error)
	DeleteMultiple(identifier string) ([]*api.Token, error)
}
