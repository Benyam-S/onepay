package service

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/wallet"
)

// Service is a type that defines user wallet service
type Service struct {
	walletRepo wallet.IWalletRepository
}

// NewWalletService is a function that returns a new user wallet service
func NewWalletService(walletRepository wallet.IWalletRepository) wallet.IService {
	return &Service{walletRepo: walletRepository}
}

// AddWallet is a method that adds a new user wallet to the system
func (service *Service) AddWallet(newWallet *entity.UserWallet) error {
	return service.walletRepo.Create(newWallet)
}

// FindWallet is a method that finds a user's wallet using the provided identifier
func (service *Service) FindWallet(identifier string) (*entity.UserWallet, error) {
	return service.walletRepo.Find(identifier)
}

// UpdateWallet is a method that updates a certain user's wallet
func (service *Service) UpdateWallet(wallet *entity.UserWallet) error {
	return service.walletRepo.Update(wallet)
}

// DeleteWallet is a method that deletes an user's wallet from the system
func (service *Service) DeleteWallet(identifier string) (*entity.UserWallet, error) {
	return service.walletRepo.Delete(identifier)
}
