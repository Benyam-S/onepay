package deleted

import "github.com/Benyam-S/onepay/entity"

// IDeletedUserRepository is an interface that defines all the repository methods for managing deleted users
type IDeletedUserRepository interface {
	Create(deletedOPUser *entity.DeletedUser) error
	Find(identifier string) (*entity.DeletedUser, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedUser
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedUser
	All(pageNum int64) []*entity.DeletedUser
	Update(deletedOPUser *entity.DeletedUser) error
	Delete(identifier string) (*entity.DeletedUser, error)
}

// IDeletedLinkedAccountRepository is an interface that defines all the repository methods for managing deleted linked accounts
type IDeletedLinkedAccountRepository interface {
	Create(deletedLinkedAccount *entity.DeletedLinkedAccount) error
	Find(identifier string) (*entity.DeletedLinkedAccount, error)
	Search(columnName string, columnValue interface{}) []*entity.DeletedLinkedAccount
	SearchMultiple(key string, pageNum int64, columns ...string) ([]*entity.DeletedLinkedAccount, int64)
	Update(deletedLinkedAccount *entity.DeletedLinkedAccount) error
	Delete(identifier string) (*entity.DeletedLinkedAccount, error)
	DeleteMultiple(identifier string) ([]*entity.DeletedLinkedAccount, error)
	IsUnique(columnName string, columnValue interface{}) bool
}

// IFrozenUserRepository is an interface that defines all the repository methods for managing frozen users
type IFrozenUserRepository interface {
	Create(frozenOPUser *entity.FrozenUser) error
	Find(identifier string) (*entity.FrozenUser, error)
	Search(key string, pageNum int64, columns ...string) []*entity.FrozenUser
	All(pageNum int64) []*entity.FrozenUser
	Update(frozenOPUser *entity.FrozenUser) error
	Delete(identifier string) (*entity.FrozenUser, error)
}

// IFrozenClientRepository is an interface that defines all the repository methods for managing frozen api clients
type IFrozenClientRepository interface {
	Create(frozenAPIClient *entity.FrozenClient) error
	Find(identifier string) (*entity.FrozenClient, error)
	Search(key string, pageNum int64, columns ...string) []*entity.FrozenClient
	All(pageNum int64) []*entity.FrozenClient
	Update(frozenAPIClient *entity.FrozenClient) error
	Delete(identifier string) (*entity.FrozenClient, error)
}
