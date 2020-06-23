package user

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a user struct
type IService interface {
	AddUser(opUser *entity.User, opPassword *entity.UserPassword) error
	ValidateUserProfile(opUser *entity.User) entity.ErrMap
	VerifyUserPassword(opPassword *entity.UserPassword, verifyPassword string) error
}
