package notifier

import (
	"bytes"
	"net/http"
)

// Notifier is a type that defines a change notifier struct
type Notifier struct {
	ListenerURI string
}

// NewNotifier is a function that returns a new notifier type
func NewNotifier(url string) *Notifier {
	return &Notifier{ListenerURI: url}
}

// NotifyWalletChange is a method that notify a certain user wallet change to its listener
func (notifier Notifier) NotifyWalletChange(id string) error {

	client := new(http.Client)
	output := bytes.NewBufferString(id)
	url := notifier.ListenerURI + "/api/v1/listener/wallet"

	request, err := http.NewRequest("PUT", url, output)
	if err != nil {
		return err
	}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
