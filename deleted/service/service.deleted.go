package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

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
	deletedLinkedAccount.AccountProviderID = linkedAccount.AccountProviderID
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

// AddStaffToTrash is a method that adds a staff member to deleted table
func (service *Service) AddStaffToTrash(staffMember *entity.Staff) error {

	deletedStaffMember := new(entity.DeletedUser)
	deletedStaffMember.UserID = staffMember.UserID
	deletedStaffMember.FirstName = staffMember.FirstName
	deletedStaffMember.LastName = staffMember.LastName
	deletedStaffMember.Email = staffMember.Email
	deletedStaffMember.PhoneNumber = staffMember.PhoneNumber

	err := service.deletedUserRepo.Create(deletedStaffMember)
	if err != nil {
		return errors.New("unable to add staff member to trash")
	}
	return nil
}

// FindDeletedUser is a method that find and return a deleted user that matchs the identifier value
func (service *Service) FindDeletedUser(identifier string) (*entity.DeletedUser, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	deletedUser, err := service.deletedUserRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no deleted user found")
	}
	return deletedUser, nil
}

// FindDeletedLinkedAccount is a method that find and return a deleted linked account that matchs the identifier value
func (service *Service) FindDeletedLinkedAccount(identifier string) (*entity.DeletedLinkedAccount, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	deletedLinkedAccount, err := service.deletedLinkedAccountRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no deleted linked account found")
	}
	return deletedLinkedAccount, nil
}

// SearchDeletedUsers is a method that searchs and returns a set of deleted users related to the key identifier
func (service *Service) SearchDeletedUsers(key, pagination string, extra ...string) []*entity.DeletedUser {

	defaultSearchColumnsRegx := []string{}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"user_id", "phone_number"}
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	result1 := make([]*entity.DeletedUser, 0)
	result2 := make([]*entity.DeletedUser, 0)
	results := make([]*entity.DeletedUser, 0)
	resultsMap := make(map[string]*entity.DeletedUser)

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results
	}

	result1 = service.deletedUserRepo.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2 = service.deletedUserRepo.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, deletedStaffMember := range result1 {
		resultsMap[deletedStaffMember.UserID] = deletedStaffMember
	}

	for _, deletedStaffMember := range result2 {
		resultsMap[deletedStaffMember.UserID] = deletedStaffMember
	}

	for _, uniqueDeletedStaffMember := range resultsMap {
		results = append(results, uniqueDeletedStaffMember)
	}

	return results
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
		linkedAccount.AccountProviderID = deletedLinkedAccount.AccountProviderID
		linkedAccount.UserID = deletedLinkedAccount.UserID
		linkedAccount.ID = deletedLinkedAccount.ID
		linkedAccounts = append(linkedAccounts, linkedAccount)
	}

	return linkedAccounts
}

// SearchMultipleDeletedLinkedAccounts is a method that searchs and returns a set of deleted linked accounts related to the key identifier
func (service *Service) SearchMultipleDeletedLinkedAccounts(key string, pageNum int64, columns ...string) ([]*entity.DeletedLinkedAccount, int64) {

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return []*entity.DeletedLinkedAccount{}, 0
	}

	return service.deletedLinkedAccountRepo.SearchMultiple(key, pageNum, columns...)
}

// FreezeUser is a method that freezs a certain user account
func (service *Service) FreezeUser(userID, reason string) error {

	empty, _ := regexp.MatchString(`^\s*$`, reason)
	if empty {
		return errors.New("reason must be provided")
	}

	frozenOPUser := new(entity.FrozenUser)
	frozenOPUser.UserID = userID
	frozenOPUser.Reason = reason

	err := service.frozenUserRepo.Create(frozenOPUser)
	if err != nil {
		return errors.New("unable to freeze account")
	}
	return nil
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

// FreezeClient is a method that freezs a certain api client
func (service *Service) FreezeClient(apiKey, reason string) error {

	empty, _ := regexp.MatchString(`^\s*$`, reason)
	if empty {
		return errors.New("reason must be provided")
	}

	FrozenClient := new(entity.FrozenClient)
	FrozenClient.APIKey = apiKey
	FrozenClient.Reason = reason

	err := service.frozenClientRepo.Create(FrozenClient)
	if err != nil {
		return errors.New("unable to freeze api client")
	}
	return nil
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
