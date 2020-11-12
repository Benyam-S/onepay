package entity

// RoleStaff is a constant that defines a staff role for a staff member table
const RoleStaff = "Staff"

// RoleAdmin is a constant that defines a admin role for a staff member table
const RoleAdmin = "Admin"

// ClientTypeWeb is a constant for web client type
const ClientTypeWeb = "Web"

// APIClientTypeInternal is a constant for internal type api client
const APIClientTypeInternal = "Internal"

// APIClientTypeExternal is a constant for external or third party client
const APIClientTypeExternal = "Third Party"

// APIClientTypeUnfiltered is a constant for unfiltered client
const APIClientTypeUnfiltered = "Unfiltered"

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

// ScopeAll is a constant that holds all usable scope values
const ScopeAll = "profile, session, send, receive, pay, wallet, history, linkedaccount, moneytoken"

// PasswordFault is a constant that holds the value password_fault-
const PasswordFault = "password_fault-"

// ReceiveFault is a constant that holds the value receive_fault-
const ReceiveFault = "receive_fault-"

// MessageIDPrefix is a constant that holds the value message_id-
const MessageIDPrefix = "message_id-"

// MessageTypeSMS is a constant that defines a message type sms
const MessageTypeSMS = "sms"

// MessageTypeEmail is a constant that defines a message type email
const MessageTypeEmail = "email"

// MessageOTPSMS is a constant that defines a message tempalate path for otp message sent through sms
const MessageOTPSMS = "/message.sms.otp.json"

// MessageVerificationSMS is a constant that defines a message tempalate path for verification message sent through sms
const MessageVerificationSMS = "/message.sms.verification.json"

// MessageVerificationEmail is a constant that defines a message tempalate path for verification message sent through email
const MessageVerificationEmail = "/message.email.verification.json"

// MessageResetEmail is a constant that defines a message tempalate path for resetting password message sent through email
const MessageResetEmail = "/message.email.reset.json"

// MessageResetSMS is a constant that defines a message tempalate path for resetting password message sent through sms
const MessageResetSMS = "/message.sms.reset.json"
