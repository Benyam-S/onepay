package user

import (
	"net/http"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/session"
)

// IService is an interface that defines all the service methods of a user struct
type IService interface {
	AddUser(opUser *entity.User, opPassword *entity.UserPassword) error
	FindUser(identifier string) (*entity.User, error)
	ValidateUserProfile(opUser *entity.User) entity.ErrMap

	VerifyUserPassword(opPassword *entity.UserPassword, verifyPassword string) error
	FindPassword(identifier string) (*entity.UserPassword, error)

	AddSession(opClientSession *session.ClientSession, opUser *entity.User, r *http.Request) error
	FindSession(identifier string) ([]*session.ServerSession, error)
	UpdateSession(opServerSession *session.ServerSession) error
	DeleteSession(identifier string) (*session.ServerSession, error)
}
