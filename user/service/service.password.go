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
	opPassword, err := service.passwordRepo.Find(identifier)
	return opPassword, err
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
		return errors.New("password does not match")
	}

	opPassword.Salt = tools.GenerateRandomString(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(opPassword.Password+opPassword.Salt), 12)
	opPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return nil
}

// UpdatePassword is a method that updates a certain's user password
func (service *Service) UpdatePassword(opPassword *entity.UserPassword) error {
	return service.passwordRepo.Update(opPassword)
}

// DeletePassword is a method that deletes a certain's user password
func (service *Service) DeletePassword(identifier string) (*entity.UserPassword, error) {
	return service.passwordRepo.Delete(identifier)
}
