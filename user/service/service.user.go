package service

import (
	"encoding/base64"
	"errors"
	"regexp"

	"github.com/Benyam-S/onepay/tools"
	"golang.org/x/crypto/bcrypt"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/user"
)

// Service is a type that defines user service
type Service struct {
	userRepo    user.IUserRepository
	passwordRep user.IPasswordRepository
}

// NewUserService is a function that returns a new user service
func NewUserService(userRepository user.IUserRepository,
	passwordRepository user.IPasswordRepository) user.IService {
	return &Service{userRepo: userRepository, passwordRep: passwordRepository}
}

// ValidateUserProfile is a method that validate a user profile.
// It checks if the user has a valid entries or not and return map of errors if any
func (service *Service) ValidateUserProfile(opUser *entity.User) map[string]error {

	errMap := make(map[string]error)
	matchFirstName, _ := regexp.MatchString(`^[a-zA-Z]\w*$`, opUser.FirstName)
	matchLastName, _ := regexp.MatchString(`^\w*$`, opUser.LastName)
	matchEmail, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, opUser.Email)
	matchPhoneNumber, _ := regexp.MatchString(`^\+\d{12}$`, opUser.PhoneNumber)

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
		errMap["phone_number"] = errors.New("phonenumber should only contain numbers and must be length of 9 digits")
	}

	if !service.userRepo.IsUnique("email", opUser.Email) {
		errMap["email"] = errors.New("email address already exists")
	}

	if !service.userRepo.IsUnique("phone_number", opUser.PhoneNumber) {
		errMap["phone_number"] = errors.New("phonenumber already exists")
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// VerifyUserPassword is a method that verify a user has provided a valid password with a matching verifypassword entry
func (service *Service) VerifyUserPassword(opPassword *entity.UserPassword, verifyPassword string) (*entity.UserPassword, error) {
	matchPassword, _ := regexp.MatchString(`^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$`, opPassword.Password)

	if len(opPassword.Password) < 8 {
		return nil, errors.New("password should contain at least 8 characters")
	}

	if !matchPassword {
		return nil, errors.New("invalid characters used in password")
	}

	if opPassword.Password != verifyPassword {
		return nil, errors.New("password does not match")
	}

	opPassword.Salt = tools.RandomStringGN(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(opPassword.Password+opPassword.Salt), 12)
	opPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return opPassword, nil
}
