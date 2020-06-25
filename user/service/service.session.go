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
