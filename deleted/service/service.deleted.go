package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// Service is a struct that defines a service that manages deleted structs
type Service struct {
	deletedUserRepo          deleted.IDeletedUserRepository
	deletedLinkedAccountRepo deleted.IDeletedLinkedAccountRepository
	frozenUserRepo           deleted.IFrozenUserRepository
	frozenClientRepo         deleted.IFrozenClientRepository
}

// NewDeletedService is a function that returns a new deleted service
func NewDeletedService(deletedUserRepository deleted.IDeletedUserRepository,
	deletedLinkedAccountRepository deleted.IDeletedLinkedAccountRepository,
	frozenUserRepository deleted.IFrozenUserRepository,
	frozenClientRepository deleted.IFrozenClientRepository) deleted.IService {

	return &Service{
		deletedUserRepo: deletedUserRepository, deletedLinkedAccountRepo: deletedLinkedAccountRepository,
		frozenUserRepo: frozenUserRepository, frozenClientRepo: frozenClientRepository}
}

// AddUserToTrash is a method that adds a onepay user to deleted table
func (service *Service) AddUserToTrash(opUser *entity.User) error {

	deletedOPUser := new(entity.DeletedUser)
	deletedOPUser.UserID = opUser.UserID
	deletedOPUser.FirstName = opUser.FirstName
	deletedOPUser.LastName = opUser.LastName
	deletedOPUser.Email = opUser.Email
	deletedOPUser.PhoneNumber = opUser.PhoneNumber

	err := service.deletedUserRepo.Create(deletedOPUser)
	if err != nil {
		return errors.New("unable to add user to trash")
	}
	return nil
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

	err := service.deletedLinkedAccountRepo.Create(deletedLinkedAccount)
	if err != nil {
		return errors.New("unable to add linked account to trash")
	}

	return nil
}

// SearchDeletedLinkedAccounts is a method that returns all the deleted linked accounts that match the given identifier
func (service *Service) SearchDeletedLinkedAccounts(columnName, columnValue string) []*entity.LinkedAccount {

	empty, _ := regexp.MatchString(`^\s*$`, columnValue)
	if empty {
		return []*entity.LinkedAccount{}
	}

	deletedLinkedAccounts := service.deletedLinkedAccountRepo.Search(columnName, columnValue)
	linkedAccounts := make([]*entity.LinkedAccount, 0)

	for _, deletedLinkedAccount := range deletedLinkedAccounts {
		linkedAccount := new(entity.LinkedAccount)
		linkedAccount.AccessToken = deletedLinkedAccount.AccessToken
		linkedAccount.AccountID = deletedLinkedAccount.AccountID
		linkedAccount.AccountProvider = deletedLinkedAccount.AccountProvider
		linkedAccount.UserID = deletedLinkedAccount.UserID
		linkedAccount.ID = deletedLinkedAccount.ID
		linkedAccounts = append(linkedAccounts, linkedAccount)
	}

	return linkedAccounts
}

// UserIsFrozen is a method that checks if a given user is frozen or not
func (service *Service) UserIsFrozen(userID string) bool {

	_, err := service.frozenUserRepo.Find(userID)
	if err != nil {
		return false
	}
	return true
}

// UnfreezeUser is a method that unfreezs a certain user account
func (service *Service) UnfreezeUser(userID string) (*entity.FrozenUser, error) {

	frozenOPUser, err := service.frozenUserRepo.Delete(userID)
	if err != nil {
		return nil, errors.New("unable to unfreeze account")
	}
	return frozenOPUser, nil
}

// ClientIsFrozen is a method that checks if a given api client is frozen or not
func (service *Service) ClientIsFrozen(apiKey string) bool {

	_, err := service.frozenClientRepo.Find(apiKey)
	if err != nil {
		return false
	}
	return true
}

// UnfreezeClient is a method that unfreezs a certain api client
func (service *Service) UnfreezeClient(apiKey string) (*entity.FrozenClient, error) {

	frozenClient, err := service.frozenClientRepo.Delete(apiKey)
	if err != nil {
		return nil, errors.New("unable to unfreeze api client")
	}
	return frozenClient, nil
}