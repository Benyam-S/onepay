package moneytoken

import "github.com/Benyam-S/onepay/entity"

// IMoneyTokenRepository is an interface that defines all the repository methods of a money token struct
type IMoneyTokenRepository interface {
	Create(newMoneyToken *entity.MoneyToken) error
	Find(identifier string) (*entity.MoneyToken, error)
	Search(identifier string) []*entity.MoneyToken
	Update(moneyToken *entity.MoneyToken) error
	UpdateValue(moneyToken *entity.MoneyToken, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.MoneyToken, error)
	DeleteMultiple(identifier string) ([]*entity.MoneyToken, error)
	IsUnique(columnName string, columnValue interface{}) bool
}
