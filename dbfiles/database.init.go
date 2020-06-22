package dbfiles

// User is a type that defines a OnePay user
type User struct {
	UserID      string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

// UserPassword is a type that defines a OnePay user password
type UserPassword struct {
	UserID   string
	Password string
	Salt     string
}
