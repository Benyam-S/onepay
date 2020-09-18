package wallet

import "github.com/Benyam-S/onepay/entity"

// IWalletRepository is an interface that defines all the repository method of a user's wallet
type IWalletRepository interface {
	Create(newOPWallet *entity.UserWallet) error
	Find(identifier string) (*entity.UserWallet, error)
	Update(opWallet *entity.UserWallet) error
	UpdateSeen(opWallet *entity.UserWallet, value bool) error
	Delete(identifier string) (*entity.UserWallet, error)
}
