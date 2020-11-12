package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Benyam-S/onepay/linkedaccount"
	"github.com/Benyam-S/onepay/moneytoken"
	"github.com/Benyam-S/onepay/notifier"
	"github.com/Benyam-S/onepay/wallet"
	"github.com/nyaruka/phonenumbers"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
)

// Service is a type that defines user service
type Service struct {
	userRepo          user.IUserRepository
	passwordRepo      user.IPasswordRepository
	preferenceRepo    user.IPreferenceRepository
	sessionRepo       user.ISessionRepository
	linkedAccountRepo linkedaccount.ILinkedAccountRepository
	moneyTokenRepo    moneytoken.IMoneyTokenRepository
	walletRepo        wallet.IWalletRepository
	apiClientRepo     user.IAPIClientRepository
	apiTokenRepo      user.IAPITokenRepository
	notifier          *notifier.Notifier
}

// NewUserService is a function that returns a new user service
func NewUserService(userRepository user.IUserRepository,
	passwordRepository user.IPasswordRepository, preferenceRepository user.IPreferenceRepository,
	sessionRepository user.ISessionRepository, apiClientRepository user.IAPIClientRepository,
	apiTokenRepository user.IAPITokenRepository, profileChangeNotifier *notifier.Notifier) user.IService {
	return &Service{userRepo: userRepository, passwordRepo: passwordRepository, preferenceRepo: preferenceRepository,
		sessionRepo: sessionRepository, apiClientRepo: apiClientRepository,
		apiTokenRepo: apiTokenRepository, notifier: profileChangeNotifier}
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

	// Since user preference is initiated with a default value it can be created here.
	userPreference := new(entity.UserPreference)
	userPreference.UserID = opUser.UserID
	err = service.preferenceRepo.Create(userPreference)
	if err != nil {
		// Cleaning up if password is not add to the database
		service.passwordRepo.Delete(opUser.UserID)
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
	validFirstName, _ := regexp.MatchString(`^[a-zA-Z]\w*$`, opUser.FirstName)
	validLastName, _ := regexp.MatchString(`^\w*$`, opUser.LastName)
	validEmail, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, opUser.Email)

	countryCode := tools.GetCountryCode(opUser.PhoneNumber)
	phoneNumber := tools.OnlyPhoneNumber(opUser.PhoneNumber)

	// Checking for local phone number
	isLocalPhoneNumber, _ := regexp.MatchString(`^0\d{9}$`, phoneNumber)
	if isLocalPhoneNumber && (countryCode == "" || countryCode == "ET") {
		phoneNumberSlice := strings.Split(phoneNumber, "")
		if phoneNumberSlice[0] == "0" {
			phoneNumberSlice = phoneNumberSlice[1:]
			internationalPhoneNumber := "+251" + strings.Join(phoneNumberSlice, "")
			phoneNumber = internationalPhoneNumber
			countryCode = "ET"
		}
	}

	parsedPhoneNumber, _ := phonenumbers.Parse(phoneNumber, "")
	validPhoneNumber := phonenumbers.IsValidNumber(parsedPhoneNumber)

	if !validFirstName {
		errMap["first_name"] = errors.New("firstname should only contain alpha numerical values and have at least one character")
	}
	if !validLastName {
		errMap["last_name"] = errors.New("lastname should only contain alpha numerical values")
	}
	if !validEmail {
		errMap["email"] = errors.New("invalid email address used")
	}

	if !validPhoneNumber {
		errMap["phone_number"] = errors.New("invalid phonenumber used")
	} else {
		// If a valid phone number is provided, adjust the phone number to fit the database
		// Stored in +251900010197[ET] or +251900010197 format
		phoneNumber = fmt.Sprintf("+%d%d", parsedPhoneNumber.GetCountryCode(),
			parsedPhoneNumber.GetNationalNumber())

		opUser.PhoneNumber = phoneNumber
		if countryCode != "" {
			opUser.PhoneNumber = fmt.Sprintf("%s[%s]", phoneNumber, countryCode)
		}
	}

	// Meaning a new user is being add
	if opUser.UserID == "" {
		if validEmail && !service.userRepo.IsUnique("email", opUser.Email) {
			errMap["email"] = errors.New("email address already exists")
		}

		phoneNumberPattern := `^` + tools.EscapeRegexpForDatabase(phoneNumber) + `(\\[[a-zA-Z]{2}])?$`
		if validPhoneNumber && !service.userRepo.IsUniqueRegx("phone_number", phoneNumberPattern) {
			errMap["phone_number"] = errors.New("phone number already exists")
		}
	} else {
		// Meaning trying to update user
		prevProfile, _ := service.userRepo.Find(opUser.UserID)

		// checking uniqueness only for email that isn't identical to the user's previous email
		if validEmail && prevProfile.Email != opUser.Email {
			if !service.userRepo.IsUnique("email", opUser.Email) {
				errMap["email"] = errors.New("email address already exists")
			}
		}

		if validPhoneNumber &&
			tools.OnlyPhoneNumber(prevProfile.PhoneNumber) != tools.OnlyPhoneNumber(opUser.PhoneNumber) {
			phoneNumberPattern := `^` + tools.EscapeRegexpForDatabase(phoneNumber) + `(\\[[a-zA-Z]{2}])?$`
			if !service.userRepo.IsUniqueRegx("phone_number", phoneNumberPattern) {
				errMap["phone_number"] = errors.New("phone number already exists")
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
		return nil, errors.New("no user found")
	}

	opUser, err := service.userRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no user found")
	}
	return opUser, nil
}

// FindUserAlsoWPhone is a method that find and return a user that matchs the identifier value,
// Also it uses the identifier as a phone number
func (service *Service) FindUserAlsoWPhone(identifier string, lb *entity.LocalizationBag) (*entity.User, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no user found")
	}

	phoneNumber := tools.OnlyPhoneNumber(identifier)
	if lb.PhoneCode != "" && !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + lb.PhoneCode + phoneNumber
		parsedPhoneNumber, _ := phonenumbers.Parse(phoneNumber, "")
		validPhoneNumber := phonenumbers.IsValidNumber(parsedPhoneNumber)
		if validPhoneNumber {
			phoneNumber = fmt.Sprintf("+%d%d", parsedPhoneNumber.GetCountryCode(),
				parsedPhoneNumber.GetNationalNumber())
		} else {
			phoneNumber = tools.OnlyPhoneNumber(identifier)
		}
	} else {
		// Localizing phone number
		splitIdentifier := strings.Split(identifier, "")
		if splitIdentifier[0] == "0" && len(splitIdentifier) == 10 {
			phoneNumber = "+251" + strings.Join(splitIdentifier[1:], "")
		}
	}

	opUser, err := service.userRepo.FindAlsoWPhone(identifier, phoneNumber)
	if err != nil {
		return nil, errors.New("no user found")
	}
	return opUser, nil
}

// AllUsers is a method that returns all the users with pagination
func (service *Service) AllUsers(pagination string) []*entity.User {
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)
	return service.userRepo.All(pageNum)
}

// SearchUsers is a method that searchs and returns a set of users related to the key identifier
func (service *Service) SearchUsers(key, pagination string, extra ...string) []*entity.User {

	defaultSearchColumnsRegx := []string{}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"user_id", "phone_number"}
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	result1 := make([]*entity.User, 0)
	result2 := make([]*entity.User, 0)
	results := make([]*entity.User, 0)
	resultsMap := make(map[string]*entity.User)

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results
	}

	result1 = service.userRepo.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2 = service.userRepo.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, opUser := range result1 {
		resultsMap[opUser.UserID] = opUser
	}

	for _, opUser := range result2 {
		resultsMap[opUser.UserID] = opUser
	}

	for _, uniqueOPUser := range resultsMap {
		results = append(results, uniqueOPUser)
	}

	return results
}

// UpdateUser is a method that updates a user in the system
func (service *Service) UpdateUser(opUser *entity.User) error {

	err := service.userRepo.Update(opUser)
	if err != nil {
		return errors.New("unable to update user")
	}

	/* ++++++++++++++ NOTIFYING CHANGE +++++++++++++++ */
	service.notifier.NotifyProfileChange(opUser.UserID)
	/* +++++++++++++++++++++++++++++++++++++++++++++++ */

	return nil
}

// UpdateUserSingleValue is a method that updates a single column entry of a user
func (service *Service) UpdateUserSingleValue(userID, columnName string, columnValue interface{}) error {

	User := entity.User{UserID: userID}
	err := service.userRepo.UpdateValue(&User, columnName, columnValue)
	if err != nil {
		return errors.New("unable to update user")
	}

	/* ++++++++++++++ NOTIFYING CHANGE +++++++++++++++ */
	service.notifier.NotifyProfileChange(userID)
	/* +++++++++++++++++++++++++++++++++++++++++++++++ */

	return nil
}

// DeleteUser is a method that deletes a user from the system including it's session's and password and other data
func (service *Service) DeleteUser(userID string) (*entity.User, error) {

	// Removing the linked tables first
	service.apiClientRepo.DeleteMultiple(userID)
	service.apiTokenRepo.DeleteMultiple(userID)
	service.passwordRepo.Delete(userID)
	service.preferenceRepo.Delete(userID)
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
