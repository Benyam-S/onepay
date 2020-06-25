package user

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/session"
)

// IUserRepository is an interface that defines all repository methods of a user struct
type IUserRepository interface {
	Create(newOPUser *entity.User) error
	Find(identifier string) (*entity.User, error)
	Update(opUser *entity.User) error
	Delete(identifier string) (*entity.User, error)
	CountUser() int
	IsUnique(columnName string, columnValue interface{}) bool
}

// IPasswordRepository is an interface that defines all repository methods of a user's password struct
type IPasswordRepository interface {
	Create(newOPPassword *entity.UserPassword) error
	Find(identifier string) (*entity.UserPassword, error)
	Update(opPassword *entity.UserPassword) error
	Delete(identifier string) (*entity.UserPassword, error)
}

// ISessionRepository is an interface that defines all repository methods of a user's server side session struct
type ISessionRepository interface {
	Create(newOPSession *session.ServerSession) error
	Find(identifier string) ([]*session.ServerSession, error)
	Update(opSession *session.ServerSession) error
	Delete(identifier string) (*session.ServerSession, error)
}
