package moneytoken

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a money token struct
type IService interface {
	AddMoneyToken(newMoneyToken *entity.MoneyToken) error
	FindMoneyToken(identifier string) (*entity.MoneyToken, error)
	SearchMoneyToken(identifier string) []*entity.MoneyToken
	ExpiredMoneyTokens() []*entity.MoneyToken
	UpdateMoneyToken(moneyToken *entity.MoneyToken) error
	UpdateMoneyTokenSingleValue(code, columnName string, columnValue interface{}) error
	DeleteMoneyToken(code string) (*entity.MoneyToken, error)
	DeleteMoneyTokens(senderID string) ([]*entity.MoneyToken, error)
}
