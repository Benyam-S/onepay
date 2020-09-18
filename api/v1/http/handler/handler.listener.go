package handler

import (
	"io/ioutil"
	"net/http"
)

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
