package entity

// WalletCheckpointError is a constant for holding the error value 'wallet checkpoint error'
const WalletCheckpointError = "wallet checkpoint error"

// MoneyTokenCheckpointError is a constant for holding the error value 'money token checkpoint error'
const MoneyTokenCheckpointError = "money token checkpoint error"

// HistoryCheckpointError is a constant for holding the error value 'history checkpoint error'
const HistoryCheckpointError = "history checkpoint error"

// TransactionBaseLimitError is a constant that holds transaction base limit error
const TransactionBaseLimitError = "amount is less than transaction base limit"

// DailyTransactionLimitError is a constant that holds daily transaction limit error
const DailyTransactionLimitError = "user has exceeded daily transaction limit"

// InsufficientBalanceError is a constant that holds insufficient balance error
const InsufficientBalanceError = "insufficient balance, please recharge your wallet"

// SenderNotFoundError is a constant that holds sender not found error
const SenderNotFoundError = "no onepay user for the provided sender id"

// ReceiverNotFoundError is a constant that holds receiver not found error
const ReceiverNotFoundError = "no onepay user for the provided receiver id"

// TransactionWSelfError is a constant that holds transaction with our own account is not allowed error
const TransactionWSelfError = "cannot make transaction with your own account"

// FrozenAccountError is a constant that holds account has been frozen error
const FrozenAccountError = "account has been frozen"

// FrozenAPIClientError is a constant that holds api client has been frozen error
const FrozenAPIClientError = "api client has been frozen"

// AmountParsingError is a constant that holds amount parsing error
const AmountParsingError = "amount parsing error"

// TooManyAttemptsError is a constant that holds too many attempts error
const TooManyAttemptsError = "too many attempts try after 24 hours"

// InvalidPasswordOrIdentifierError is a constant that holds invalid password or identifier error
const InvalidPasswordOrIdentifierError = "invalid identifier or password used"

// InvalidPasswordError is a constant that holds invalid password used error
const InvalidPasswordError = "invalid password used"

// InternalAPIClientError is a constant that holds unable to add an internal api client error
const InternalAPIClientError = "unable to add an internal api client"

// APITokenError is a constant that holds unable to create an api token error
const APITokenError = "unable to create an api token"
