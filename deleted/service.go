package deleted

import "github.com/Benyam-S/onepay/entity"

// IService is a method that defines all the service methods for managing deleted structs
type IService interface {
	AddUserToTrash(opUser *entity.User) error
	AddLinkedAccountToTrash(linkedAccount *entity.LinkedAccount) error
	SearchDeletedLinkedAccounts(columnName, columnValue string) []*entity.LinkedAccount

	UserIsFrozen(userID string) bool
	UnfreezeUser(userID string) (*entity.FrozenUser, error)
	ClientIsFrozen(apiKey string) bool
	UnfreezeClient(apiKey string) (*entity.FrozenClient, error)
}
