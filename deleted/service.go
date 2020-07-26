package deleted

import "github.com/Benyam-S/onepay/entity"

// IService is a method that defines all the service methods for managing deleted structs
type IService interface {
	AddUserToTrash(opUser *entity.User) error
	AddStaffToTrash(staffMember *entity.Staff) error
	AddLinkedAccountToTrash(linkedAccount *entity.LinkedAccount) error

	FindDeletedUser(identifier string) (*entity.DeletedUser, error)
	SearchDeletedUsers(key, pagination string, extra ...string) []*entity.DeletedUser
	FindDeletedLinkedAccount(identifier string) (*entity.DeletedLinkedAccount, error)
	SearchDeletedLinkedAccounts(columnName, columnValue string) []*entity.LinkedAccount
	SearchMultipleDeletedLinkedAccounts(key, pagination string, columns ...string) []*entity.DeletedLinkedAccount

	FreezeUser(userID, reason string) error
	UserIsFrozen(userID string) bool
	UnfreezeUser(userID string) (*entity.FrozenUser, error)
	FreezeClient(apiKey, reason string) error
	ClientIsFrozen(apiKey string) bool
	UnfreezeClient(apiKey string) (*entity.FrozenClient, error)
}
