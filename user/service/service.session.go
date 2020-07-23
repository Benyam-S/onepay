package service

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/Benyam-S/onepay/client/http/session"
	"github.com/Benyam-S/onepay/entity"
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
		return errors.New("unable to add new session")
	}
	return nil
}

// FindSession is a method that finds and return a user's server side session that matchs the identifier value
func (service *Service) FindSession(identifier string) (*session.ServerSession, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	serverSession, err := service.sessionRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("session not found")
	}
	return serverSession, nil
}

// SearchSession is a method that searchs and return a user's server side session that matchs the identifier value
func (service *Service) SearchSession(identifier string) ([]*session.ServerSession, error) {

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("empty identifier used")
	}

	serverSessions, err := service.sessionRepo.Search(identifier)
	if err != nil {
		return nil, errors.New("no session found for the provided identifier")
	}
	return serverSessions, nil
}

// UpdateSession is a method that updates a user's server side session
func (service *Service) UpdateSession(opServerSession *session.ServerSession) error {

	err := service.sessionRepo.Update(opServerSession)
	if err != nil {
		return errors.New("unable to update session")
	}
	return nil
}

// DeleteSession is a method that deletes a user's server side session from the system
func (service *Service) DeleteSession(identifier string) (*session.ServerSession, error) {

	opServerSession, err := service.sessionRepo.Delete(identifier)
	if err != nil {
		return nil, errors.New("unable to delete session")
	}
	return opServerSession, nil
}
