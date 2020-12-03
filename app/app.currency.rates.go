package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
)

// GetCurrencyRates is a method that retrive currency rates/exchange values for the provided base currency
func (onepay *OnePay) GetCurrencyRates(base string) ([]*entity.CurrencyRate, error) {

	if !tools.IsValidBaseCurrency(base) {
		return nil, errors.New("unknown base currency")
	}

	listOfCurrencyRates := make([]*entity.CurrencyRate, 0)
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd,
		fmt.Sprintf("./assets/currency/currency.rates.%s.json", strings.ToUpper(base)))

	output, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(output, &listOfCurrencyRates)
	if err != nil {
		return nil, err
	}

	return listOfCurrencyRates, nil
}
