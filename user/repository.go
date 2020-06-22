package user

import "github.com/Benyam-S/onepay/entity"

// IUserRepository is an interface that defines all repository methods of a user struct
type IUserRepository interface {
	Create(newOPUser *entity.User) (*entity.User, error)
	Find(identifier string) (*entity.User, error)
	Update(opUser *entity.User) (*entity.User, error)
	Delete(identifier string) (*entity.User, error)
	CountUser() int
	IsUnique(columnName string, columnValue interface{}) bool
}

// IPasswordRepository is an interface that defines all repository methods of a user's password struct
type IPasswordRepository interface {
	Create(newOPPassword *entity.UserPassword) (*entity.UserPassword, error)
	Find(identifier string) (*entity.UserPassword, error)
	Update(opPassword *entity.UserPassword) (*entity.UserPassword, error)
	Delete(identifier string) (*entity.UserPassword, error)
}
