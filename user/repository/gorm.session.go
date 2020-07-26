package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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
// In Find() session_id is only used as a key
func (repo *SessionRepository) Find(identifier string) (*session.ServerSession, error) {

	opSession := new(session.ServerSession)
	err := repo.conn.Model(opSession).
		Where("session_id = ?", identifier).
		Find(opSession).Error

	if err != nil {
		return nil, err
	}

	return opSession, nil
}

// Search is a method that searchs for a set of server side sessions from the database using an identifier.
// In Search() user_id is only used as a key
func (repo *SessionRepository) Search(identifier string) ([]*session.ServerSession, error) {

	var opSessions []*session.ServerSession
	err := repo.conn.Model(session.ServerSession{}).
		Where("user_id = ?", identifier).
		Find(&opSessions).Error

	if err != nil {
		return nil, err
	}

	if len(opSessions) == 0 {
		return nil, errors.New("no available session for the provided identifier")
	}
	return opSessions, nil
}

// SearchMultiple is a method that search and returns a set of server side sessions from that matchs the key identifier.
func (repo *SessionRepository) SearchMultiple(key string, pageNum int64, columns ...string) []*session.ServerSession {

	var opSessions []*session.ServerSession
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM server_sessions WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY user_id ASC LIMIT ?, 30", sqlValues...).Scan(&opSessions)

	return opSessions
}

// SearchMultipleWRegx is a method that searchs and returns set of server side sessions limited to the key identifier and page number using regular experssions
func (repo *SessionRepository) SearchMultipleWRegx(key string, pageNum int64, columns ...string) []*session.ServerSession {

	var opSessions []*session.ServerSession
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM server_sessions WHERE "+strings.Join(whereStmt, "||")+" ORDER BY user_id ASC LIMIT ?, 30", sqlValues...).Scan(&opSessions)

	return opSessions
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
// In Delete() session_id is only used as a key
func (repo *SessionRepository) Delete(identifier string) (*session.ServerSession, error) {
	opSession := new(session.ServerSession)
	err := repo.conn.Model(session.ServerSession{}).Where("session_id = ?", identifier).First(opSession).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opSession)
	return opSession, nil
}

// DeleteMultiple is a method that deletes multiple user's server side session from the database use identifier.
// In DeleteMultiple() user_id is only used as a key
func (repo *SessionRepository) DeleteMultiple(identifier string) ([]*session.ServerSession, error) {
	var userSessions []*session.ServerSession
	err := repo.conn.Model(session.ServerSession{}).Where("user_id = ?", identifier).Find(&userSessions).Error

	if err != nil {
		return nil, err
	}

	if len(userSessions) == 0 {
		return nil, errors.New("no session for the provided identifier")
	}

	repo.conn.Model(session.ServerSession{}).Where("user_id = ?", identifier).Delete(session.ServerSession{})
	return userSessions, nil
}
