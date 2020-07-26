package tools

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

// GenerateRandomString is a function that generate a random string based on a given length.
func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomBytes returns securely generated random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateOTP is a function that generates a random otp value of 4 digits
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	nBig := rand.Int63n(8999)
	return fmt.Sprintf("%d", nBig+1000)
}

// GenerateMoneyTokenCode is a function that generates a code for money token struct
func GenerateMoneyTokenCode() string {

	rand.Seed(time.Now().Unix())
	gnFourDigits := func() string {
		num1 := rand.Int63n(9)
		num2 := rand.Int63n(9)
		num3 := rand.Int63n(9)
		num4 := rand.Int63n(9)

		if num1 == num2 && num1 == num3 && num1 == num4 {
			return ""
		}

		output := fmt.Sprintf("%d%d%d%d", num1, num2, num3, num4)
		return output
	}

	code := ""

	for i := 0; i < 4; i++ {
		part := gnFourDigits()
		if part != "" {
			code += part
		} else {
			i--
		}
	}

	return code

}

// IDWOutPrefix is a function that returns an id without it's prefix
func IDWOutPrefix(id string) string {

	var output string
	prefixes := []string{`OP_API-`, `OP_Token-`, `OP_LA-`, `deleted-\w{4}:`, `OP_S-`, `OP-`}

	for _, prefix := range prefixes {

		match, _ := regexp.MatchString(`^`+prefix, regexp.QuoteMeta(id))
		if match {

			rx := regexp.MustCompile(prefix)
			output = rx.ReplaceAllString(id, "")

			break
		}

	}

	return output
}
