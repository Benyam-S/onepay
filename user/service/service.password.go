package service

import (
	"encoding/base64"
	"errors"
	"regexp"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"golang.org/x/crypto/bcrypt"
)

// FindPassword is a method that find and return a user's password that matchs the identifier value
func (service *Service) FindPassword(identifier string) (*entity.UserPassword, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("password not found")
	}

	opPassword, err := service.passwordRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("password not found")
	}

	return opPassword, nil
}

// VerifyUserPassword is a method that verify a user has provided a valid password with a matching verifypassword entry
func (service *Service) VerifyUserPassword(opPassword *entity.UserPassword, verifyPassword string) error {
	matchPassword, _ := regexp.MatchString(`^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$`, opPassword.Password)

	if len(opPassword.Password) < 8 {
		return errors.New("password should contain at least 8 characters")
	}

	if !matchPassword {
		return errors.New("invalid characters used in password")
	}

	if opPassword.Password != verifyPassword {
		return errors.New("passwords do not match")
	}

	opPassword.Salt = tools.GenerateRandomString(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(opPassword.Password+opPassword.Salt), 12)
	opPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return nil
}

// UpdatePassword is a method that updates a certain user's password
func (service *Service) UpdatePassword(opPassword *entity.UserPassword) error {

	err := service.passwordRepo.Update(opPassword)
	if err != nil {
		return errors.New("unable to update password")
	}
	return nil
}

// DeletePassword is a method that deletes a certain user's password
func (service *Service) DeletePassword(identifier string) (*entity.UserPassword, error) {

	opPassword, err := service.passwordRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete password")
	}
	return opPassword, nil
}
