package wallet

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a user wallet struct
type IService interface {
	AddWallet(newWallet *entity.UserWallet) error
	FindWallet(identifier string) (*entity.UserWallet, error)
	UpdateWallet(wallet *entity.UserWallet) error
	DeleteWallet(identifier string) (*entity.UserWallet, error)
}
