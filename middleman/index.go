package middleman

import "github.com/Benyam-S/onepay/entity"

// GetAccountInfo is
func GetAccountInfo(accountID, accessToken string) (*entity.AccountInfo, error) {
	return &entity.AccountInfo{Amount: 200}, nil
}

// RefillAccount is
func RefillAccount(accountID, accessToken string, amount float64) error {
	return nil
}

// WithdrawFromAccount is
func WithdrawFromAccount(accountID, accessToken string, amount float64) error {
	return nil
}

// AddLinkedAccount is
func AddLinkedAccount(accountID, accountProvider string) error {
	return nil
}

// VerifyLinkedAccount is
func VerifyLinkedAccount(otp string) (string, error) {
	return "1234589", nil
}
