package tools

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Benyam-S/onepay/entity"
)

// SendSMS is a function that sends a given message to the provide phone number
func SendSMS(to, msg string) (string, error) {

	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "./assets/accounts", "/account.api.sms.json")
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return "", err
	}

	var clientAccount entity.APIClientSMS
	err = json.Unmarshal(data, &clientAccount)
	if err != nil {
		return "", err
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + clientAccount.AccountID + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", clientAccount.From)
	msgData.Set("Body", msg)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(clientAccount.AccountID, clientAccount.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			return "", err
		}

		smsID, ok := data["sid"].(string)
		if !ok {
			return "", errors.New("unable to parse the sms id")
		}
		return smsID, nil
	}

	return "", errors.New(resp.Status)
}
