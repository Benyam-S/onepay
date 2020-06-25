package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
)

// Service is a type that defines user service
type Service struct {
	userRepo     user.IUserRepository
	passwordRepo user.IPasswordRepository
	sessionRepo  user.ISessionRepository
}

// NewUserService is a function that returns a new user service
func NewUserService(userRepository user.IUserRepository,
	passwordRepository user.IPasswordRepository, sessionRepository user.ISessionRepository) user.IService {
	return &Service{userRepo: userRepository, passwordRepo: passwordRepository, sessionRepo: sessionRepository}
}

// AddUser is a method that adds a new OnePay user to the system along with the password
func (service *Service) AddUser(opUser *entity.User, opPassword *entity.UserPassword) error {
	err := service.userRepo.Create(opUser)
	if err != nil {
		return err
	}
	opPassword.UserID = opUser.UserID
	err = service.passwordRepo.Create(opPassword)
	if err != nil {
		// Cleaning up if password is not add to the database
		service.userRepo.Delete(opUser.UserID)
		return err
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

	if matchEmail && !service.userRepo.IsUnique("email", opUser.Email) {
		errMap["email"] = errors.New("email address already exists")
	}

	if matchPhoneNumber && !service.userRepo.IsUnique("phone_number", opUser.PhoneNumber) {
		errMap["phone_number"] = errors.New("phonenumber already exists")
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindUser is a method that find and return a user that matchs the identifier value
func (service *Service) FindUser(identifier string) (*entity.User, error) {
	opUser, err := service.userRepo.Find(identifier)
	return opUser, err
}
