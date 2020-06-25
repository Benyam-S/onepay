package service

import (
	"net/http"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/session"
)

// AddSession is a method that adds a new user session to the system using the client side session
func (service *Service) AddSession(opClientSession *session.ClientSession, opUser *entity.User, r *http.Request) error {

	opServerSession := new(session.ServerSession)
	opServerSession.SessionID = opClientSession.SessionID
	opServerSession.UserID = opUser.UserID
	opServerSession.DeviceInfo = r.UserAgent()
	opServerSession.IPAddress = r.Host

	err := service.sessionRepo.Create(opServerSession)
	if err != nil {
		return err
	}
	return nil
}

// FindSession is a method that finds and return a user's server side session that matchs the identifier value
func (service *Service) FindSession(identifier string) ([]*session.ServerSession, error) {
	opSession, err := service.sessionRepo.Find(identifier)
	return opSession, err
}

// UpdateSession is a method that updates a user's server side session
func (service *Service) UpdateSession(opServerSession *session.ServerSession) error {
	err := service.sessionRepo.Update(opServerSession)
	return err
}

// DeleteSession is a method that deletes a user's server side session from the system
func (service *Service) DeleteSession(identifier string) (*session.ServerSession, error) {
	opServerSession, err := service.sessionRepo.Delete(identifier)
	return opServerSession, err
}
