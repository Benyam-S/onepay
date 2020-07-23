package handler

// ErrorBody is a simple struct for holding errors
type ErrorBody struct {
	Error string `xml:"error" json:"error"`
}

// CodeBody is a simple struct for holding money token struct
type CodeBody struct {
	Code string `xml:"code" json:"code"`
}

// LinkedAccountBody is a struct that contain both the account info the linked account values
type LinkedAccountBody struct {
	ID              string
	AccountID       string
	AccountProvider string
	Amount          float64
}
