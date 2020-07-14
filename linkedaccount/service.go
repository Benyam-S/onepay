package linkedaccount

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a linked account struct
type IService interface {
	AddLinkedAccount(newLinkedAccount *entity.LinkedAccount) error
	VerifyToLink(linkedAccount *entity.LinkedAccount) error
	FindLinkedAccount(identifier string) (*entity.LinkedAccount, error)
	SearchLinkedAccounts(columnName string, columnValue interface{}) []*entity.LinkedAccount
	UpdateLinkedAccount(linkedAccount *entity.LinkedAccount) error
	UpdateLinkedAccountSingleValue(id, columnName string, columnValue interface{}) error
	DeleteLinkedAccount(id string) (*entity.LinkedAccount, error)
	DeleteLinkedAccounts(userID string) ([]*entity.LinkedAccount, error)
}
