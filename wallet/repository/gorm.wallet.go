package repository

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/wallet"
	"github.com/jinzhu/gorm"
)

// WalletRepository is a type that defines a user's wallet repository
type WalletRepository struct {
	conn *gorm.DB
}

// NewWalletRepository is a function that returns a new user's wallet repository
func NewWalletRepository(connection *gorm.DB) wallet.IWalletRepository {
	return &WalletRepository{conn: connection}
}

// Create is a method that adds a new user wallet to the database
func (repo *WalletRepository) Create(newOPWallet *entity.UserWallet) error {

	err := repo.conn.Create(newOPWallet).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain user's wallet from the database using an identifier.
// In Find() user_id is only used as a key
func (repo *WalletRepository) Find(identifier string) (*entity.UserWallet, error) {
	opWallet := new(entity.UserWallet)
	err := repo.conn.Model(opWallet).
		Where("user_id = ?", identifier).
		First(opWallet).Error

	if err != nil {
		return nil, err
	}
	return opWallet, nil
}

// Update is a method that updates a certain user's wallet value in the database
func (repo *WalletRepository) Update(opWallet *entity.UserWallet) error {

	prevOPWallet := new(entity.UserWallet)
	err := repo.conn.Model(prevOPWallet).Where("user_id = ?", opWallet.UserID).First(prevOPWallet).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(opWallet).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateSeen is a method that updates a certain user wallet's seen value in the database
func (repo *WalletRepository) UpdateSeen(opWallet *entity.UserWallet, value bool) error {

	prevOPWallet := new(entity.UserWallet)
	err := repo.conn.Model(prevOPWallet).Where("user_id = ?", opWallet.UserID).First(prevOPWallet).Error

	if err != nil {
		return err
	}

	// Since we are only changing the seen value we have to keep the previous updated time
	err = repo.conn.Exec("UPDATE user_wallets SET seen = ? WHERE user_id = ?", value, opWallet.UserID).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user's wallet from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *WalletRepository) Delete(identifier string) (*entity.UserWallet, error) {
	opWallet := new(entity.UserWallet)
	err := repo.conn.Model(opWallet).Where("user_id = ?", identifier).First(opWallet).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opWallet)
	return opWallet, nil
}
