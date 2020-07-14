package service

import (
	"fmt"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// Service is a struct that defines a service that manages deleted structs
type Service struct {
	deletedUserRepo          deleted.IDeletedUserRepository
	deletedLinkedAccountRepo deleted.IDeletedLinkedAccountRepository
}

// NewDeletedService is a function that returns a new deleted service
func NewDeletedService(deletedUserRepository deleted.IDeletedUserRepository,
	deletedLinkedAccountRepository deleted.IDeletedLinkedAccountRepository) deleted.IService {
	return &Service{deletedUserRepo: deletedUserRepository,
		deletedLinkedAccountRepo: deletedLinkedAccountRepository}
}

// AddUserToTrash is a method that adds a onepay user to deleted table
func (service *Service) AddUserToTrash(opUser *entity.User) error {

	deletedOPUser := new(entity.DeletedUser)
	deletedOPUser.UserID = opUser.UserID
	deletedOPUser.FirstName = opUser.FirstName
	deletedOPUser.LastName = opUser.LastName
	deletedOPUser.Email = opUser.Email
	deletedOPUser.PhoneNumber = opUser.PhoneNumber

	return service.deletedUserRepo.Create(deletedOPUser)
}

// AddLinkedAccountToTrash is a method that adds a linked account to deleted table
func (service *Service) AddLinkedAccountToTrash(linkedAccount *entity.LinkedAccount) error {

	deletedLinkedAccount := new(entity.DeletedLinkedAccount)
	deletedLinkedAccount.AccessToken = linkedAccount.AccessToken
	deletedLinkedAccount.AccountID = linkedAccount.AccountID
	deletedLinkedAccount.AccountProvider = linkedAccount.AccountProvider
	deletedLinkedAccount.UserID = linkedAccount.UserID
	deletedLinkedAccount.ID = fmt.Sprintf("deleted-%s:", tools.GenerateRandomString(4)) + linkedAccount.ID

	// Creating a unique id for the deleted linked account by adding 'deleted-xxxx'
	for !service.deletedLinkedAccountRepo.IsUnique("id", deletedLinkedAccount.ID) {
		deletedLinkedAccount.ID = fmt.Sprintf("deleted-%s:", tools.GenerateRandomString(4)) + linkedAccount.ID
	}

	return service.deletedLinkedAccountRepo.Create(deletedLinkedAccount)
}
