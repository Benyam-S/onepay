package dbfiles

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// Init is a function that will initialize the OnePay database
func Init(db *gorm.DB) {

	db.AutoMigrate(&entity.UserPassword{})
	db.AutoMigrate(&entity.User{})
	db.Model(&entity.UserPassword{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
}
