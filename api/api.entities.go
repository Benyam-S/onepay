package api

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Client is a type that defines a OnePay api client
type Client struct {
	ClientUserID string `gorm:"not null"`
	APIKey       string `gorm:"unique; not null"`
	APISecret    string `gorm:"not null"`
	CallBack     string `gorm:"not null"`
	APPName      string
	Type         string `gorm:"not null"`
	gorm.Model
}

// Token is a type that defines a OnePay api access token
type Token struct {
	AccessToken     string `gorm:"not null; unique"`
	UserID          string `gorm:"not null"`
	APIKey          string `gorm:"not null"`
	ExpiresAt       int64
	DailyExpiration int64
	gorm.Model
}

// TableName is a method that set Tokens's table name to be `api_tokens`
func (Token) TableName() string {
	return "api_tokens"
}

// TableName is a method that set Clients's table name to be `api_clients`
func (Client) TableName() string {
	return "api_clients"
}

// Valid a is a method that ensures Token is type jwt.Clamis
func (apiToken Token) Valid() error {
	if time.Now().Unix() > apiToken.ExpiresAt {
		return errors.New("invalid token, api token has expired")
	}

	return nil
}

// PastDailyExpiration is a method that checks if the api token has exceeded daily expiration time
func (apiToken Token) PastDailyExpiration() bool {

	if time.Now().Unix() > time.Unix(apiToken.DailyExpiration, 0).Add(time.Hour*24).Unix() {
		return true
	}

	return false
}

// Extract is a function that generate a valid api token from a signed string
func Extract(signedToken string) (*Token, error) {

	token, err := jwt.ParseWithClaims(signedToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error in signing method")
		}
		return []byte(os.Getenv("onepay_secret_key")), nil
	})

	if err != nil {
		return nil, err
	}

	apiToken, ok := token.Claims.(*Token)
	if !ok || apiToken.Valid() != nil {
		return nil, errors.New("invalid api token")
	}

	return apiToken, nil

}
