package entity

// User is a type that defines a OnePay user
type User struct {
	UserID      string `gorm:"primary_key; unique; not null"`
	FirstName   string `gorm:"not null"`
	LastName    string
	Email       string `gorm:"not null; unique"`
	PhoneNumber string `gorm:"not null; unique"`
	ProfilePic  string
}

// UserPassword is a type that defines a OnePay user password
type UserPassword struct {
	UserID   string `gorm:"primary_key; unique; not null"`
	Password string `gorm:"not null"`
	Salt     string `gorm:"not null"`
}
