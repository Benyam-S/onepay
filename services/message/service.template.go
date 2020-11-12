package message

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Benyam-S/onepay/entity"
)

// CreateMessageBodyFromTemplate is a function that creates a message from pre defined templates
func CreateMessageBodyFromTemplate(path string, inputs ...string) (string, error) {

	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "./assets/messages", path)
	data, err := ioutil.ReadFile(dir)

	if err != nil {
		return "", err
	}

	var messageTemplate map[string][]string
	err = json.Unmarshal(data, &messageTemplate)

	if err != nil {
		return "", err
	}

	var body string

	switch path {
	case entity.MessageOTPSMS:
		body = messageTemplate["message_body"][0] + inputs[0] + ". " + messageTemplate["message_body"][1]
	case entity.MessageVerificationEmail:
		body = messageTemplate["message_body"][0] + inputs[0] +
			messageTemplate["message_body"][1] + inputs[1] +
			". " + messageTemplate["message_body"][2]
	case entity.MessageVerificationSMS:
		body = messageTemplate["message_body"][0] + inputs[0] +
			messageTemplate["message_body"][1] + inputs[1] +
			". " + messageTemplate["message_body"][2]
	case entity.MessageResetEmail:
		fallthrough
	case entity.MessageResetSMS:
		body = messageTemplate["message_body"][0] + inputs[0] + ". " + messageTemplate["message_body"][1]
	}

	return body, nil
}
