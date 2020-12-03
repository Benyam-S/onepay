package tools

import "strings"

// IsValidBaseCurrency is a function that check whether the provided base currency is valid or not
func IsValidBaseCurrency(base string) bool {

	listOfValidBaseCurrencies := []string{"ETB"}
	isValid := false

	for _, validBaseCurrency := range listOfValidBaseCurrencies {
		if validBaseCurrency == strings.ToUpper(base) {
			isValid = true
			break
		}
	}

	return isValid
}
