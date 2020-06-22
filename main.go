package main

import (
	"fmt"

	"github.com/Benyam-S/onepay/dbfiles"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "root:0911@tcp(127.0.0.1:3306)/onepay?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to the database: mysql @GORM")

	dbfiles.Init(db)

	// userRepo := repository.NewUserRepository(db)
	// passwordRepo := repository.NewPasswordRepository(db)

	// tempUser := entity.User{UserID: "OPbB49Kkw2", FirstName: "Benyam", LastName: "Simayehu"}
	// tempPassword := entity.UserPassword{UserID: "OPh7lTo5t1", Password: "12443", Salt: "123"}

	// user, err := userRepo.Delete("OPh7lTo5t1")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(user)

	// password, err := passwordRepo.Delete("OPh7lTo5t1")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(password)
}
