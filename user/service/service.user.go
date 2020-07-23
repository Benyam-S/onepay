package service

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Benyam-S/onepay/linkedaccount"
	"github.com/Benyam-S/onepay/moneytoken"
	"github.com/Benyam-S/onepay/wallet"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
)

// Service is a type that defines user service
type Service struct {
	userRepo          user.IUserRepository
	passwordRepo      user.IPasswordRepository
	sessionRepo       user.ISessionRepository
	linkedAccountRepo linkedaccount.ILinkedAccountRepository
	moneyTokenRepo    moneytoken.IMoneyTokenRepository
	walletRepo        wallet.IWalletRepository
	apiClientRepo     user.IAPIClientRepository
	apiTokenRepo      user.IAPITokenRepository
}

// NewUserService is a function that returns a new user service
func NewUserService(userRepository user.IUserRepository,
	passwordRepository user.IPasswordRepository, sessionRepository user.ISessionRepository,
	apiClientRepository user.IAPIClientRepository, apiTokenRepository user.IAPITokenRepository) user.IService {
	return &Service{userRepo: userRepository, passwordRepo: passwordRepository, sessionRepo: sessionRepository,
		apiClientRepo: apiClientRepository, apiTokenRepo: apiTokenRepository}
}

// AddUser is a method that adds a new OnePay user to the system along with the password
func (service *Service) AddUser(opUser *entity.User, opPassword *entity.UserPassword) error {
	err := service.userRepo.Create(opUser)
	if err != nil {
		return errors.New("unable to add new user")
	}
	opPassword.UserID = opUser.UserID
	err = service.passwordRepo.Create(opPassword)
	if err != nil {
		// Cleaning up if password is not add to the database
		service.userRepo.Delete(opUser.UserID)
		return errors.New("unable to add new user")
	}

	return nil
}

// ValidateUserProfile is a method that validate a user profile.
// It checks if the user has a valid entries or not and return map of errors if any.
// Also it will add country code to the phone number value if not included: default country code +251
func (service *Service) ValidateUserProfile(opUser *entity.User) entity.ErrMap {

	errMap := make(map[string]error)
	matchFirstName, _ := regexp.MatchString(`^[a-zA-Z]\w*$`, opUser.FirstName)
	matchLastName, _ := regexp.MatchString(`^\w*$`, opUser.LastName)
	matchEmail, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, opUser.Email)
	matchPhoneNumber, _ := regexp.MatchString(`^(\+\d{11,12})|(0\d{9})$`, opUser.PhoneNumber)

	if !matchFirstName {
		errMap["first_name"] = errors.New("firstname should only contain alpha numerical values and have at least one character")
	}
	if !matchLastName {
		errMap["last_name"] = errors.New("lastname should only contain alpha numerical values")
	}
	if !matchEmail {
		errMap["email"] = errors.New("invalid email address used")
	}

	if !matchPhoneNumber {
		errMap["phone_number"] = errors.New("phonenumber should be +XXXXXXXXXXXX or 0XXXXXXXXX formate, also use url escaping if country code was used")
	} else {
		// If a valid phone number is provided, adjust the phone number to fit the database
		phoneNumberSlice := strings.Split(opUser.PhoneNumber, "")
		if phoneNumberSlice[0] == "0" {
			phoneNumberSlice = phoneNumberSlice[1:]
			validPhoneNumber := "+251" + strings.Join(phoneNumberSlice, "")
			opUser.PhoneNumber = validPhoneNumber
		}
	}

	// Meaning a new user is being add
	if opUser.UserID == "" {
		if matchEmail && !service.userRepo.IsUnique("email", opUser.Email) {
			errMap["email"] = errors.New("email address already exists")
		}

		if matchPhoneNumber && !service.userRepo.IsUnique("phone_number", opUser.PhoneNumber) {
			errMap["phone_number"] = errors.New("phonenumber already exists")
		}
	} else {
		// Meaning trying to update user
		prevProfile, _ := service.userRepo.Find(opUser.UserID)

		// checking uniquness only for email that isn't identical to the user's previous email
		if matchEmail && prevProfile.Email != opUser.Email {
			if !service.userRepo.IsUnique("email", opUser.Email) {
				errMap["email"] = errors.New("email address already exists")
			}
		}

		if matchPhoneNumber && prevProfile.PhoneNumber != opUser.PhoneNumber {
			if !service.userRepo.IsUnique("phone_number", opUser.PhoneNumber) {
				errMap["phone_number"] = errors.New("phonenumber already exists")
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindUser is a method that find and return a user that matchs the identifier value
func (service *Service) FindUser(identifier string) (*entity.User, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	opUser, err := service.userRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no user found")
	}
	return opUser, nil
}

// UpdateUser is a method that updates a user in the system
func (service *Service) UpdateUser(opUser *entity.User) error {

	err := service.userRepo.Update(opUser)
	if err != nil {
		return errors.New("unable to update user")
	}
	return nil
}

// UpdateUserSingleValue is a method that updates a single column entiry of a user
func (service *Service) UpdateUserSingleValue(userID, columnName string, columnValue interface{}) error {

	User := entity.User{UserID: userID}
	err := service.userRepo.UpdateValue(&User, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update user")
	}
	return nil
}

// DeleteUser is a method that deletes a user from the system including it's session's and password and other datas
func (service *Service) DeleteUser(userID string) (*entity.User, error) {

	// Removing the linked tables first
	service.apiClientRepo.DeleteMultiple(userID)
	service.apiTokenRepo.DeleteMultiple(userID)
	service.passwordRepo.Delete(userID)
	service.sessionRepo.DeleteMultiple(userID)

	opUser, err := service.userRepo.Delete(userID)
	if err != nil {
		return nil, errors.New("unable to delete user")
	}

	if opUser.ProfilePic != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "./assets/profilepics", opUser.ProfilePic)
		tools.RemoveFile(filePath)
	}

	return opUser, nil
}
