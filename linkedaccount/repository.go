package linkedaccount

import "github.com/Benyam-S/onepay/entity"

// ILinkedAccountRepository is an inteface that defines all the repository methods of a linked account struct
type ILinkedAccountRepository interface {
	Create(newLinkedAccount *entity.LinkedAccount) error
	Find(identifier string) (*entity.LinkedAccount, error)
	Search(colunmName string, columnValue interface{}) []*entity.LinkedAccount
	Update(linkedAccount *entity.LinkedAccount) error
	UpdateValue(linkedAccount *entity.LinkedAccount, columnName string, columnValue interface{}) error
	Delete(identifier string) (*entity.LinkedAccount, error)
	DeleteMultiple(identifier string) ([]*entity.LinkedAccount, error)
	IsUnique(columnName string, columnValue interface{}) bool
}
