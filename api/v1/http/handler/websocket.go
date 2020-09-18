package handler

import (
	"net/http"
	"time"

	"github.com/Benyam-S/onepay/tools"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/Benyam-S/onepay/entity"
)

// HandleCreateWebsocket is a handler func that creates a websocket connection with client
func (handler *UserAPIHandler) HandleCreateWebsocket(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	opUser, ok := ctx.Value(entity.Key("onepay_user")).(*entity.User)

	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	format := mux.Vars(r)["format"]

	handler.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	websocketChannel := make(chan interface{})
	closeChannel := make(chan bool)
	handler.pushWebsocketChannel(websocketChannel, opUser.UserID)

	go func() {
		for {
			time.Sleep(time.Second * 60)
			err = ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second*50))
			if err != nil {
				closeChannel <- true
				return
			}
		}
	}()

	for {

		select {

		case change := <-websocketChannel:
			output, _ := tools.MarshalIndent(change, "", "\t", format)
			ws.WriteMessage(websocket.TextMessage, output)

		case <-closeChannel:
			ws.Close()
			handler.removeWebsocketChannel(websocketChannel, opUser.UserID)
			close(websocketChannel)
			close(closeChannel)
			return
		}
	}

}

func (handler *UserAPIHandler) pushWebsocketChannel(channel chan interface{}, key string) {
	handler.Lock()
	defer handler.Unlock()

	if handler.activeSocketChannels == nil {
		handler.activeSocketChannels = make(map[string][]chan interface{})
	}

	prevChannels := handler.activeSocketChannels[key]
	handler.activeSocketChannels[key] = append(prevChannels, channel)

}

func (handler *UserAPIHandler) removeWebsocketChannel(channel chan interface{}, key string) {
	handler.Lock()
	defer handler.Unlock()

	prevChannels := handler.activeSocketChannels[key]
	newChannels := make([]chan interface{}, 0)

	for _, ch := range prevChannels {
		if ch == channel {
			continue
		}
		newChannels = append(newChannels, ch)
	}

	handler.activeSocketChannels[key] = newChannels
}
