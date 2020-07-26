package app

import "github.com/Benyam-S/onepay/entity"

// UserHistory is a method that returns the user history only
func (onepay *OnePay) UserHistory(userID string, pagenation int64, viewBys ...string) []*entity.UserHistory {

	orderBy := "id"
	opHistories := make([]*entity.UserHistory, 0)
	length := len(viewBys)

	for _, viewBy := range viewBys {

		searchColumns := make([]string, 0)
		methods := make([]string, 0)

		if viewBy == "transaction_received" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode)

		} else if viewBy == "transaction_sent" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode)

		} else if viewBy == "payment_received" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodPaymentQRCode)

		} else if viewBy == "payment_sent" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodPaymentQRCode)

		} else if viewBy == "recharged" {
			if length == 1 {
				orderBy = "received_at"
			}
			searchColumns = append(searchColumns, "receiver_id")
			methods = append(methods, entity.MethodRecharged)

		} else if viewBy == "withdrawn" {
			if length == 1 {
				orderBy = "sent_at"
			}
			searchColumns = append(searchColumns, "sender_id")
			methods = append(methods, entity.MethodWithdrawn)

		} else if viewBy == "all" && length == 1 {
			searchColumns = append(searchColumns, "sender_id", "receiver_id")
			methods = append(methods, entity.MethodTransactionOnePayID,
				entity.MethodTransactionQRCode, entity.MethodPaymentQRCode)
		} else {
			// If it is unkown view by
			continue
		}

		result := onepay.HistoryService.SearchHistories(userID, orderBy, methods, pagenation, searchColumns...)
		opHistories = append(opHistories, result...)

	}

	return opHistories
}
