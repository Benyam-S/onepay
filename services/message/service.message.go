package message

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/go-redis/redis"
)

// StartMessageServices is a function that starts the messaging service
func StartMessageServices(redisClient *redis.Client, serviceChannel chan *entity.MessageTemp) {
	for {
		message := <-serviceChannel
		err := SendMessage(message)

		if err == nil {
			output, _ := json.MarshalIndent(message, "\t", "")
			tools.SetValue(redisClient, message.ID, string(output), time.Hour*6)
		}
	}
}

// SendMessage is a function that sends message according to its type
func SendMessage(message *entity.MessageTemp) error {
	if message.Type == entity.MessageTypeSMS {
		tools.SendSMS(tools.OnlyPhoneNumber(message.To), message.Body)
		return nil
	} else if message.Type == entity.MessageTypeEmail {
		tools.SendEmail(message.To, message.Subject, message.Body)
		return nil
	}

	return errors.New("invalid message type")
}
