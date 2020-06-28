package repository

import (
	"errors"

	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// SessionRepository is a type that defines a user server side session repository
type SessionRepository struct {
	conn *gorm.DB
}

// NewSessionRepository is a function that returns a new user server side session repository
func NewSessionRepository(connection *gorm.DB) user.ISessionRepository {
	return &SessionRepository{conn: connection}
}

// Create is a method that adds a new user session to the database
func (repo *SessionRepository) Create(newOPSession *session.ServerSession) error {

	err := repo.conn.Create(newOPSession).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that find user's server side sessions from the database using an identifier.
// In Find() user_id and session_id can be used as an key
func (repo *SessionRepository) Find(identifier string) ([]*session.ServerSession, error) {
	var opSessions []*session.ServerSession
	err := repo.conn.Model(session.ServerSession{}).
		Where("user_id = ? || session_id = ?", identifier, identifier).
		Find(&opSessions).Error

	if err != nil {
		return nil, err
	}

	if len(opSessions) == 0 {
		return nil, errors.New("no available session for the provided identifier")
	}
	return opSessions, nil
}

// Update is a method that updates a certain user's server side Session value in the database
func (repo *SessionRepository) Update(opSession *session.ServerSession) error {

	prevOPSession := new(session.ServerSession)
	err := repo.conn.Model(prevOPSession).Where("session_id = ?", opSession.SessionID).First(prevOPSession).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(session.ServerSession{}).Where("session_id = ?", opSession.SessionID).Update(opSession).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user's server side session from the database using an identifier.
// In Delete() session_id is only used as an key
func (repo *SessionRepository) Delete(identifier string) (*session.ServerSession, error) {
	opSession := new(session.ServerSession)
	err := repo.conn.Model(session.ServerSession{}).Where("session_id = ?", identifier).First(opSession).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opSession)
	return opSession, nil
}
