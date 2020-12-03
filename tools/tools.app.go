package tools

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/Benyam-S/onepay/entity"
)

// ClosingStatement is a function that generates a file that contain a user histories and linked account information
// for a user that is deleting it's onepay account
func ClosingStatement(opUser *entity.Staff, histories []*entity.UserHistory,
	linkedAccounts []*entity.LinkedAccount) (string, error) {

	fileName := opUser.UserID + "_" + GenerateRandomString(7) + ".txt"
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "./assets/statements", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString("\n\n ********************************* USER PROFILE ********************************* \n\n")
	if err != nil {
		return "", err
	}

	line := "First Name: " + opUser.FirstName + "	\nLast Name: " + opUser.LastName +
		" \nEmail: " + opUser.Email + " \nPHONE NUMBER: " + opUser.PhoneNumber

	_, err = file.WriteString(line)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString("\n\n ********************************* USER LINKED ACCOUNTS ********************************* \n\n")
	if err != nil {
		return "", err
	}

	for _, linkedAccount := range linkedAccounts {
		line := "Account ID: " + linkedAccount.AccountID + "	Account Provider: " + linkedAccount.AccountProviderID + "\n"
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
