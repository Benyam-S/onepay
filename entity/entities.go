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

// ErrMap is a type that defines a map with string identifier and it's error
type ErrMap map[string]error

// StringMap is a method that returns string map corresponding to the ErrMap where the error type is converted to a string
func (errMap ErrMap) StringMap() map[string]string {
	stringMap := make(map[string]string)
	for key, value := range errMap {
		stringMap[key] = value.Error()
	}

	return stringMap
}
