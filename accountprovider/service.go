package accountprovider

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a account provider struct
type IService interface {
	AddAccountProvider(newAccountProvider *entity.AccountProvider) error
	ValidateAccountProvider(accountProvider *entity.AccountProvider) entity.ErrMap
	FindAccountProvider(identifier string) (*entity.AccountProvider, error)
	AllAccountProviders() []*entity.AccountProvider
	SearchAccountProvider(key, orderBy string, pageNum int64, columns ...string) ([]*entity.AccountProvider, int64)
	UpdateAccountProvider(accountProvider *entity.AccountProvider) error
	UpdateAccountProviderSingleValue(id, columnName string, columnValue interface{}) error
	DeleteAccountProvider(id string) (*entity.AccountProvider, error)
}
