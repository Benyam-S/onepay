package deleted

import "github.com/Benyam-S/onepay/entity"

// IDeletedUserRepository is an inteface that defines all the repository methods for managing deleted users
type IDeletedUserRepository interface {
	Create(deletedOPUser *entity.DeletedUser) error
	Find(identifier string) (*entity.DeletedUser, error)
	Update(deletedOPUser *entity.DeletedUser) error
	Delete(identifier string) (*entity.DeletedUser, error)
}

// IDeletedLinkedAccountRepository is an inteface that defines all the repository methods for managing deleted linked accounts
type IDeletedLinkedAccountRepository interface {
	Create(deletedLinkedAccount *entity.DeletedLinkedAccount) error
	Find(identifier string) (*entity.DeletedLinkedAccount, error)
	Update(deletedLinkedAccount *entity.DeletedLinkedAccount) error
	Delete(identifier string) (*entity.DeletedLinkedAccount, error)
	DeleteMultiple(identifier string) ([]*entity.DeletedLinkedAccount, error)
	IsUnique(columnName string, columnValue interface{}) bool
}
