package deleted

import "github.com/Benyam-S/onepay/entity"

// IService is a method that defines all the service methods for managing deleted structs
type IService interface {
	AddUserToTrash(opUser *entity.User) error
	AddLinkedAccountToTrash(linkedAccount *entity.LinkedAccount) error
}
