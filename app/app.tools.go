package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/go-redis/redis"
)

// GetTransactionFee is a function that returns the appropirate transaction fee for the provided amount
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
