package handler

import "github.com/Benyam-S/onepay/entity"

// ErrorBody is a simple struct for holding errors
type ErrorBody struct {
	Error string `xml:"error" json:"error"`
}

// CodeBody is a simple struct for holding money token struct
type CodeBody struct {
	Code string `xml:"code" json:"code"`
}

// HistoriesContainer is a struct that contain a single request histories with it's page count
type HistoriesContainer struct {
	Result      []*entity.UserHistory
	CurrentPage int64
	PageCount   int64
}

// NotifierContainer is a struct that holds a change notifier value
type NotifierContainer struct {
	Type string
	Body interface{}
}
