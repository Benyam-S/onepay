package entity

import (
	"net/http"
	"time"
)

// User is a type that defines a OnePay user
type User struct {
	UserID      string `gorm:"primary_key; unique; not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null; unique"`
	PhoneNumber string `gorm:"not null; unique"`
	ProfilePic  string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Staff is a type that defines a staff member
type Staff struct {
	UserID      string `gorm:"primary_key; unique; not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	PhoneNumber string `gorm:"unique; not null"`
	Email       string `gorm:"unique; not null"`
	ProfilePic  string `gorm:"not null"`
	Role        string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserPassword is a type that defines a OnePay user password
type UserPassword struct {
	UserID   string `gorm:"primary_key; unique; not null"`
	Password string `gorm:"not null"`
	Salt     string `gorm:"not null"`
}

// UserWallet is a type that defines a OnePay user wallet
type UserWallet struct {
	UserID    string  `gorm:"primary_key; unique; not null"`
	Amount    float64 `gorm:"not null"`
	Seen      bool    `gorm:"default: true;"`
	UpdatedAt time.Time
}

// UserHistory is a type that defines a OnePay user's history
type UserHistory struct {
	ID           int    `gorm:"primary_key; unique; not null"`
	SenderID     string `gorm:"not null"`
	ReceiverID   string `gorm:"not null"`
	SentAt       time.Time
	ReceivedAt   time.Time
	Method       string  `gorm:"not null"`
	Code         string  `gorm:"not null"`
	Amount       float64 `gorm:"not null"`
	SenderSeen   bool    `gorm:"default: false;"`
	ReceiverSeen bool    `gorm:"default: false;"`
}

// UserPreference is a type that defines a OnePay user preference
type UserPreference struct {
	UserID              string `gorm:"primary_key; unique; not null"`
	TwoStepVerification bool   `gorm:"not null; default: false"`
}

// MoneyToken is a type that defines a token generated for qr code
type MoneyToken struct {
	Code           string `gorm:"primary_key; unique; not null"`
	SenderID       string `gorm:"not null"`
	SentAt         time.Time
	Amount         float64 `gorm:"not null"`
	ExpirationDate time.Time
	Method         string `gorm:"not null"`
}

// LinkedAccount is a type that defines an account that is linked with OnePay account
type LinkedAccount struct {
	ID                string `gorm:"primary_key; unique; not null"`
	UserID            string `gorm:"not null"`
	AccountProviderID string `gorm:"not null"`
	AccountID         string `gorm:"not null"`
	AccessToken       string `gorm:"not null"`
}

// AccountProvider is type that defines an external account provider
type AccountProvider struct {
	ID        string `gorm:"primary_key; unique; not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
}

// Extras is a type that defines values that are required but extra in definition
type Extras struct {
	TotalUsersCount int
}

// DeletedUser is a type that defines a OnePay user that has been deleted
// This struct is used to store and identify a pervious user
type DeletedUser struct {
	UserID      string `gorm:"primary_key; unique; not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null"`
	PhoneNumber string `gorm:"not null"`
}

// DeletedLinkedAccount is a type that defines an account that was linked with OnePay account
// It can be used which account has been linked to which
type DeletedLinkedAccount struct {
	ID                string `gorm:"primary_key; unique; not null"`
	UserID            string `gorm:"not null"`
	AccountProviderID string `gorm:"not null"`
	AccountID         string `gorm:"not null"`
	AccessToken       string `gorm:"not null"`
}

// DeletedAccountProvider is a type that defines an account provider that has been a linked from OnePay System
type DeletedAccountProvider struct {
	ID   string `gorm:"primary_key; unique; not null"`
	Name string `gorm:"not null"`
}

// FrozenUser is a struct that defines a user that has been frozen or deactivated
type FrozenUser struct {
	UserID    string `gorm:"primary_key; unique; not null"`
	Reason    string `gorm:"not null"`
	CreatedAt time.Time
}

// FrozenClient is a struct that defines an api client that has been frozen or deactivated
type FrozenClient struct {
	APIKey    string `gorm:"primary_key; unique; not null"`
	Reason    string `gorm:"not null"`
	CreatedAt time.Time
}

// TableName is a method that defines the database table name of the user history struct
func (UserHistory) TableName() string {
	return "user_history"
}

// AccountInfo is type that defines an external account information
type AccountInfo struct {
	Amount            float64
	AccountID         string
	AccountProviderID string
}

// LinkedAccountContainer is a struct that contains a filtered linked account without unnecessary values
type LinkedAccountContainer struct {
	ID                  string
	UserID              string
	AccountID           string
	AccountProviderName string
	AccountProviderID   string
}

// LocalizationBag is a struct that contains localization values
type LocalizationBag struct {
	CountryCode string
	PhoneCode   string
}

// MessageTemp is a struct that defines what the message contains and its type
type MessageTemp struct {
	ID      string
	To      string
	Type    string
	Body    string
	Subject string
}

// CurrencyRate is a struct that defines a single currency rate/exchange value
type CurrencyRate struct {
	FromSymbol   string
	FromName     string
	ToSymbol     string
	ToName       string
	CurrentValue float64
	Values       []float64
	Dates        []time.Time
}

// Key is a type that defines a key type that can be used a key value in context
type Key string

// Middleware is a type that defines a function that takes a handler func and return a new handler func type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// ErrMap is a type that defines a map with string identifier and it's error
type ErrMap map[string]error

// StringMap is a method that returns string map corresponding to the ErrMap where the error type is converted to a string
func (errMap ErrMap) StringMap() map[string]string {
	stringMap := make(map[string]string)
	for key, value := range errMap {
		stringMap[key] = value.Error()
	}

	return stringMap
}

// Equal is a method that checks if the two history objects are identical
func (history *UserHistory) Equal(opHistory *UserHistory) bool {

	check1 := history.ID == opHistory.ID
	check2 := history.Code == opHistory.Code
	check3 := history.Amount == opHistory.Amount
	check4 := history.Method == opHistory.Method
	check5 := history.ReceivedAt.Unix() == opHistory.ReceivedAt.Unix()
	check6 := history.SentAt.Unix() == opHistory.SentAt.Unix()
	check7 := history.ReceiverID == opHistory.ReceiverID
	check8 := history.SenderID == opHistory.SenderID

	return check1 && check2 && check3 && check4 && check5 && check6 && check7 && check8
}

// Equal is a method that checks if the two history objects are identical
func (wallet *UserWallet) Equal(opWallet *UserWallet) bool {

	check1 := wallet.UserID == opWallet.UserID
	check2 := wallet.Amount == opWallet.Amount

	return check1 && check2
}

// Equal is a method that checks if the two money token objects are identical
func (moneyToken *MoneyToken) Equal(opMoneyToken *MoneyToken) bool {

	check1 := moneyToken.Code == opMoneyToken.Code
	check2 := moneyToken.Amount == opMoneyToken.Amount
	check3 := moneyToken.SenderID == opMoneyToken.SenderID
	check4 := moneyToken.SentAt.Unix() == opMoneyToken.SentAt.Unix()
	check5 := moneyToken.ExpirationDate.Unix() == opMoneyToken.ExpirationDate.Unix()
	check6 := moneyToken.Method == opMoneyToken.Method

	return check1 && check2 && check3 && check4 && check5 && check6
}
