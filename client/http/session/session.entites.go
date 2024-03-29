package session

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Benyam-S/onepay/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// ClientSession is a type that defines a client side user session
type ClientSession struct {
	SessionID       string
	ExpiresAt       int64
	DailyExpiration int64
	CreatedAt       int64
	UpdatedAt       int64
}

// ServerSession is a type that defines a servide side user session
type ServerSession struct {
	SessionID   string `gorm:"primary_key; unique; not null"`
	UserID      string `gorm:"not null"`
	IPAddress   string `gorm:"not null"`
	DeviceInfo  string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deactivated bool // This can be used to identify a session that has been logged out
}

// Valid a is a method that ensures session is type jwt.Claims
func (session ClientSession) Valid() error {
	if time.Now().Unix() > session.ExpiresAt {
		return errors.New("invalid session, session has expired")
	}

	return nil
}

// PastDailyExpiration is a method that checks if the session has exceeded daily expiration time
func (session ClientSession) PastDailyExpiration() bool {

	if time.Now().Unix() > time.Unix(session.DailyExpiration, 0).Add(time.Hour*24).Unix() {
		return true
	}

	return false
}

// Save is a method that save a given session as a client cookie on the network
func (session ClientSession) Save(w http.ResponseWriter) error {
	signedString, err := tools.GenerateToken([]byte(os.Getenv("onepay_secret_key")), session)
	if err != nil {
		return err
	}

	maxAge := session.ExpiresAt - time.Now().Unix()

	clientCookie := http.Cookie{
		Name:     os.Getenv("onepay_cookie_name"),
		Value:    signedString,
		MaxAge:   int(maxAge),
		Expires:  time.Unix(maxAge, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, &clientCookie)
	return nil
}

// Remove expires existing session
func (session ClientSession) Remove(w http.ResponseWriter) {
	c := http.Cookie{
		Name:    os.Getenv("onepay_cookie_name"),
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Value:   "",
	}
	http.SetCookie(w, &c)
}

// Create is a function that creates a new session type with the given userID and current time
func Create(userID string) *ClientSession {

	newOpSession := new(ClientSession)
	newOpSession.CreatedAt = time.Now().Unix()
	newOpSession.DailyExpiration = time.Now().Unix()
	newOpSession.UpdatedAt = time.Now().Unix()
	newOpSession.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()

	sessionUUID := uuid.Must(uuid.NewRandom())
	newOpSession.SessionID = sessionUUID.String()

	return newOpSession
}

// Extract is a function that generate a valid session from a signed string
func Extract(signedToken string) (*ClientSession, error) {

	token, err := jwt.ParseWithClaims(signedToken, &ClientSession{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error in signing method")
		}
		return []byte(os.Getenv("onepay_secret_key")), nil
	})

	if err != nil {
		return nil, err
	}

	opSession, ok := token.Claims.(*ClientSession)
	if !ok || opSession.Valid() != nil {
		return nil, errors.New("invalid session")
	}

	return opSession, nil

}
