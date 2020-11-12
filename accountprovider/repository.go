package accountprovider

import "github.com/Benyam-S/onepay/entity"

// IAccountProviderRepository is an interface that defines all the repository methods of account provider struct
type IAccountProviderRepository interface {
	Create(newAccountProvider *entity.AccountProvider) error
	Find(identifier string) (*entity.AccountProvider, error)
	All() []*entity.AccountProvider
	Search(columnName string, columnValue interface{}) []*entity.AccountProvider
	SearchMultiple(key, orderBy string, pageNum int64, columns ...string) ([]*entity.AccountProvider, int64)
	Update(accountProvider *entity.AccountProvider) error
	UpdateValue(accountProvider *entity.AccountProvider, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.AccountProvider, error)
	IsUnique(columnName string, columnValue interface{}) bool
}
