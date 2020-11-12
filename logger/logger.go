package logger

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/Benyam-S/onepay/entity"
)

// Logger is a type that defines a logger that log a fatal errors that needs to be reloaded
type Logger struct {
	walletLogDir     string
	moneyTokenLogDir string
	historyLogDir    string
}

// NewLogger is a function that returns a logger for OnePay fatal errors
func NewLogger(path string) *Logger {

	walletLogDirectory := filepath.Join(path, "log.wallet.json")
	historyLogDirectory := filepath.Join(path, "log.history.json")
	moneyTokenLogDirectory := filepath.Join(path, "log.money.token.json")

	return &Logger{walletLogDir: walletLogDirectory, historyLogDir: historyLogDirectory,
		moneyTokenLogDir: moneyTokenLogDirectory}
}

// Must is a function that panics if the provided logger method returns error
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// LoggedWallets is a method that returns all the logged wallets.
// In another word all other user wallets that failed to update
func (logger *Logger) LoggedWallets() []*entity.UserWallet {

	loggedWallets := make([]*entity.UserWallet, 0)

	loggedWalletsByte, err := ioutil.ReadFile(logger.walletLogDir)
	if err == nil {
		json.Unmarshal(loggedWalletsByte, &loggedWallets)
	}

	return loggedWallets
}

// LoggedHistories is a method that returns all the logged histories.
func (logger *Logger) LoggedHistories() []*entity.UserHistory {

	loggedHistories := make([]*entity.UserHistory, 0)

	loggedHistoriesByte, err := ioutil.ReadFile(logger.historyLogDir)
	if err == nil {
		json.Unmarshal(loggedHistoriesByte, &loggedHistories)
	}

	return loggedHistories
}

// LoggedMoneyTokens is a method that returns all the logged money tokens
func (logger *Logger) LoggedMoneyTokens() []*entity.MoneyToken {

	loggedMoneyTokens := make([]*entity.MoneyToken, 0)

	loggedMoneyTokensByte, err := ioutil.ReadFile(logger.moneyTokenLogDir)
	if err == nil {
		json.Unmarshal(loggedMoneyTokensByte, &loggedMoneyTokens)
	}

	return loggedMoneyTokens
}

// LogWallet is a method that logs a wallet struct in log file
func (logger *Logger) LogWallet(opWallet *entity.UserWallet) error {

	loggedWallets := make([]*entity.UserWallet, 0)

	loggedWalletsByte, err := ioutil.ReadFile(logger.walletLogDir)
	if err == nil {
		json.Unmarshal(loggedWalletsByte, &loggedWallets)
	}

	loggedWallets = append(loggedWallets, opWallet)
	output, err := json.MarshalIndent(loggedWallets, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.walletLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}

// LogHistory is a method that logs a history struct in log file
func (logger *Logger) LogHistory(opHistory *entity.UserHistory) error {

	loggedHistories := make([]*entity.UserHistory, 0)

	loggedHistoriesByte, err := ioutil.ReadFile(logger.historyLogDir)
	if err == nil {
		json.Unmarshal(loggedHistoriesByte, &loggedHistories)
	}

	loggedHistories = append(loggedHistories, opHistory)
	output, err := json.MarshalIndent(loggedHistories, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.historyLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}

// LogMoneyToken is a method that logs a money token struct in log file
func (logger *Logger) LogMoneyToken(moneyToken *entity.MoneyToken) error {

	moneyTokens := make([]*entity.MoneyToken, 0)

	moneyTokensByte, err := ioutil.ReadFile(logger.moneyTokenLogDir)
	if err == nil {
		json.Unmarshal(moneyTokensByte, &moneyTokens)
	}

	moneyTokens = append(moneyTokens, moneyToken)
	output, err := json.MarshalIndent(moneyTokens, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.moneyTokenLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}

// RemoveWallet is a method that removes a wallet object from the log file
func (logger *Logger) RemoveWallet(opWallet *entity.UserWallet) error {

	loggedWallets := make([]*entity.UserWallet, 0)
	loggedWalletsWithout := make([]*entity.UserWallet, 0)
	addedFlag := false

	loggedWalletsByte, err := ioutil.ReadFile(logger.walletLogDir)
	if err == nil {
		json.Unmarshal(loggedWalletsByte, &loggedWallets)
	}

	for _, loggedWallet := range loggedWallets {

		if loggedWallet.Equal(opWallet) && !addedFlag {
			// This ensures to deleted only one wallet object that is identical to a wallet in the file log
			addedFlag = true
			continue
		}
		loggedWalletsWithout = append(loggedWalletsWithout, loggedWallet)
	}

	output, err := json.MarshalIndent(loggedWalletsWithout, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.walletLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}

// RemoveHistory is a method that removes a user history object from the log file
func (logger *Logger) RemoveHistory(opHistory *entity.UserHistory) error {

	loggedHistories := make([]*entity.UserHistory, 0)
	loggedHistoriesWithout := make([]*entity.UserHistory, 0)
	addedFlag := false

	loggedHistoriesByte, err := ioutil.ReadFile(logger.historyLogDir)
	if err == nil {
		json.Unmarshal(loggedHistoriesByte, &loggedHistories)
	}

	for _, loggedHistory := range loggedHistories {

		if loggedHistory.Equal(opHistory) && !addedFlag {
			// This ensures to deleted only one history object that is identical to a history in the file log
			addedFlag = true
			continue
		}
		loggedHistoriesWithout = append(loggedHistoriesWithout, loggedHistory)
	}

	output, err := json.MarshalIndent(loggedHistoriesWithout, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.historyLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}

// RemoveMoneyToken is a method that removes a money token object from the log file
func (logger *Logger) RemoveMoneyToken(opMoneyToken *entity.MoneyToken) error {

	moneyTokens := make([]*entity.MoneyToken, 0)
	moneyTokensWithout := make([]*entity.MoneyToken, 0)
	addedFlag := false

	moneyTokensByte, err := ioutil.ReadFile(logger.moneyTokenLogDir)
	if err == nil {
		json.Unmarshal(moneyTokensByte, &moneyTokens)
	}

	for _, moneyToken := range moneyTokens {

		if moneyToken.Equal(opMoneyToken) && !addedFlag {
			// This ensures to deleted only one money token object that is identical to a money token in the file log
			addedFlag = true
			continue
		}
		moneyTokensWithout = append(moneyTokensWithout, moneyToken)
	}

	output, err := json.MarshalIndent(moneyTokensWithout, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(logger.moneyTokenLogDir, output, 7777)
	if err != nil {
		return err
	}

	return nil
}
