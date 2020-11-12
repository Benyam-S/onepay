package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/go-redis/redis"
)

// GetTransactionFee is a function that returns the appropriate transaction fee for the provided amount
func GetTransactionFee(amount float64) float64 {

	transactionFee, _ := strconv.ParseFloat(os.Getenv(entity.TransactionFee), 64)
	return transactionFee
}

// AboveTransactionBaseLimit is a function that checks if the provided amount is above the transaction base limit
func AboveTransactionBaseLimit(amount float64) bool {

	transactionBaseLimit, _ := strconv.ParseFloat(os.Getenv(entity.TransactionBaseLimit), 64)
	if amount >= transactionBaseLimit {
		return true
	}
	return false
}

// AboveWithdrawBaseLimit is a function that checks if the provided amount is above the withdraw base limit
func AboveWithdrawBaseLimit(amount float64) bool {

	withdrawBaseLimit, _ := strconv.ParseFloat(os.Getenv(entity.WithdrawBaseLimit), 64)
	if amount >= withdrawBaseLimit {
		return true
	}
	return false
}

// AboveDailyTransactionLimit is a method that checks wheather a user has exceeded the daily transaction limit
func AboveDailyTransactionLimit(userID string, amount float64, redisClient *redis.Client) bool {

	prevAmountString, _ := tools.GetValue(redisClient, "daily_transaction_limit:"+userID)
	prevAmount, _ := strconv.ParseFloat(prevAmountString, 64)
	dailyTransactionLimit, _ := strconv.ParseFloat(os.Getenv(entity.DailyTransactionLimit), 64)
	if (prevAmount + amount) > dailyTransactionLimit {
		return true
	}

	return false
}

// AddToDailyTransaction is a method that adds a certain amount to a user daily transaction limit
func AddToDailyTransaction(userID string, amount float64, redisClient *redis.Client) error {

	prevAmountString, _ := tools.GetValue(redisClient, "daily_transaction_limit:"+userID)
	prevAmount, _ := strconv.ParseFloat(prevAmountString, 64)
	currentAmountString := fmt.Sprintf("%f", prevAmount+amount)

	return tools.SetValue(redisClient, "daily_transaction_limit:"+userID, currentAmountString, time.Hour*24)
}

// ClosingStatement is a function that generates a file that contain a user histories and linked account information
// for a user that is deleting it's onepay account
func ClosingStatement(opUser *entity.User, histories []*entity.UserHistory,
	linkedAccounts []*entity.LinkedAccountContainer) (string, error) {

	fileName := opUser.UserID + "_" + tools.GenerateRandomString(7) + ".txt"
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "./assets/temp", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString("\n\n ********************************* USER PROFILE ********************************* \n\n")
	if err != nil {
		return "", err
	}

	line := "First Name: " + opUser.FirstName + "	\nLast Name: " + opUser.LastName +
		" \nEmail: " + opUser.Email + " \nPHONE NUMBER: " + tools.OnlyPhoneNumber(opUser.PhoneNumber)

	_, err = file.WriteString(line)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString("\n\n ********************************* USER LINKED ACCOUNTS ********************************* \n\n")
	if err != nil {
		return "", err
	}

	for _, linkedAccount := range linkedAccounts {
		line := "Account ID: " + linkedAccount.AccountID +
			"	Account Provider: " + linkedAccount.AccountProviderName + "\n"
		_, err = file.WriteString(line)
		if err != nil {
			return "", err
		}
	}

	_, err = file.WriteString("\n\n ********************************* USER HISTORY ********************************* \n\n")
	if err != nil {
		return "", err
	}

	for _, history := range histories {
		line := "Sender ID: " + history.SenderID + "		Receiver ID: " + history.ReceiverID +
			"	Sent At: " + history.SentAt.String() + "	Received At: " + history.ReceivedAt.String() +
			"	Method: " + history.Method + "	Amount:" + strconv.FormatFloat(history.Amount, 'f', 2, 64) + "\n"

		_, err = file.WriteString(line)
		if err != nil {
			return "", err
		}
	}

	return fileName, nil
}
