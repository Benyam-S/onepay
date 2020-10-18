package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Benyam-S/onepay/entity"
)

// HandleListenToProfileChange is a handler func that listens to profile change from its notifier
func (handler *UserAPIHandler) HandleListenToProfileChange(w http.ResponseWriter, r *http.Request) {

	handler.Lock()
	defer handler.Unlock()

	userIDB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	userID := string(userIDB)
	opUser, err := handler.uService.FindUser(userID)
	if err != nil {
		return
	}

	activeSocketChannels := handler.activeSocketChannels[userID]
	for _, channel := range activeSocketChannels {
		channel <- NotifierContainer{Type: "user", Body: opUser}
	}
}

// HandleListenToWalletChange is a handler func that listens to wallet change from its notifier
func (handler *UserAPIHandler) HandleListenToWalletChange(w http.ResponseWriter, r *http.Request) {

	handler.Lock()
	defer handler.Unlock()

	userIDB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	userID := string(userIDB)
	opWallet, err := handler.app.WalletService.FindWallet(userID)
	if err != nil {
		return
	}

	activeSocketChannels := handler.activeSocketChannels[userID]
	for _, channel := range activeSocketChannels {
		channel <- NotifierContainer{Type: "wallet", Body: opWallet}
	}
}

// HandleListenToHistoryChange is a handler func that listens to history change from its notifier
func (handler *UserAPIHandler) HandleListenToHistoryChange(w http.ResponseWriter, r *http.Request) {

	handler.Lock()
	defer handler.Unlock()

	history := new(entity.UserHistory)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &history)
	if err != nil {
		return
	}

	senderActiveSocketChannels := handler.activeSocketChannels[history.SenderID]
	receiverActiveSocketChannels := handler.activeSocketChannels[history.ReceiverID]

	for _, channel := range senderActiveSocketChannels {
		channel <- NotifierContainer{Type: "history", Body: history}
	}

	for _, channel := range receiverActiveSocketChannels {
		channel <- NotifierContainer{Type: "history", Body: history}
	}
}
