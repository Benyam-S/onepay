package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Benyam-S/onepay/entity"
)

// Notifier is a type that defines a change notifier struct
type Notifier struct {
	ListenerURI string
}

// NewNotifier is a function that returns a new notifier type
func NewNotifier(url string) *Notifier {
	return &Notifier{ListenerURI: url}
}

// NotifyProfileChange is a method that notify a certain user profile change to its listener
func (notifier Notifier) NotifyProfileChange(id string) error {

	client := new(http.Client)
	output := bytes.NewBufferString(id)
	url := notifier.ListenerURI + "/api/v1/listener/profile"

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

// NotifyHistoryChange is a method that notify a certain user history change to its listener
func (notifier Notifier) NotifyHistoryChange(history *entity.UserHistory) error {

	client := new(http.Client)
	jsonOutput, _ := json.MarshalIndent(history, "", "\t")
	output := bytes.NewBuffer(jsonOutput)
	url := notifier.ListenerURI + "/api/v1/listener/history"

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
