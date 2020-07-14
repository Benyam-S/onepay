package entity

// APIClientTypeInternal is a constant for internal type api client
const APIClientTypeInternal = "Internal"

// APIClientTypeExternal is a constant for external or third party client
const APIClientTypeExternal = "Third Party"

// APIClientAppNameInternal is a constant for the app name of an internal type api client
const APIClientAppNameInternal = "OnePay"

// MethodTransactionQRCode is a constant that defines a transaction via qr code
const MethodTransactionQRCode = "Transaction Via QR Code"

// MethodPaymentQRCode is a constant that defines a payment done with qr code
const MethodPaymentQRCode = "Payment Via QR Code"

// MethodTransactionOnePayID is a constant that defines a transaction via OnePay id
const MethodTransactionOnePayID = "Transaction Via OnePay ID"

// MethodRecharged is a constant that defines an account has been recharged
const MethodRecharged = "Recharged"

// MethodWithdrawn is a constant that defines money has been withdrawn from account
const MethodWithdrawn = "Withdrawn"

// TransactionFee is a constant for holding the transaction_fee name
const TransactionFee = "transaction_fee"

// TransactionBaseLimit is a constant for holding the transaction_base_limit name
const TransactionBaseLimit = "transaction_base_limit"

// WithdrawBaseLimit is a constant for holding the withdraw_base_limit name
const WithdrawBaseLimit = "withdraw_base_limit"

// DailyTransactionLimit is a constant for holding the daily_transaction_limit name
const DailyTransactionLimit = "daily_transaction_limit"
