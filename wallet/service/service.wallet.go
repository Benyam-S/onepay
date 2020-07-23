package service

import (
	"errors"
	"regexp"

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

	err := service.walletRepo.Create(newWallet)
	if err != nil {
		return errors.New("unable to add new user wallet")
	}
	return nil
}

// FindWallet is a method that finds a user's wallet using the provided identifier
func (service *Service) FindWallet(identifier string) (*entity.UserWallet, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	opWallet, err := service.walletRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("user wallet not found")
	}
	return opWallet, nil
}

// UpdateWallet is a method that updates a certain user's wallet
func (service *Service) UpdateWallet(wallet *entity.UserWallet) error {

	err := service.walletRepo.Update(wallet)
	if err != nil {
		return errors.New("unable to update user wallet")
	}
	return nil
}

// DeleteWallet is a method that deletes an user's wallet from the system
func (service *Service) DeleteWallet(identifier string) (*entity.UserWallet, error) {

	opWallet, err := service.walletRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete user wallet")
	}
	return opWallet, nil
}
