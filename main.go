package main

import (
	"fmt"
	"net/http"

	"github.com/Benyam-S/onepay/client/http/handler"
	"github.com/Benyam-S/onepay/user/service"
	"github.com/go-redis/redis"

	"github.com/Benyam-S/onepay/dbfiles"
	"github.com/Benyam-S/onepay/user/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	db, err := gorm.Open("mysql", "root:0911@tcp(127.0.0.1:3306)/onepay?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to the database: mysql @GORM")

	dbfiles.Init(db)

	userRepo := repository.NewUserRepository(db)
	passwordRepo := repository.NewPasswordRepository(db)
	userService := service.NewUserService(userRepo, passwordRepo)

	// tempUser := entity.User{UserID: "OPMRpo8kn1", FirstName: "Benyam",
	// LastName: "Simayehu", Email: "binysimayehu@gmail.co", PhoneNumber: "+25191173268"}
	// tempPassword := entity.UserPassword{UserID: "OPh7lTo5t1", Password: "12443", Salt: "123"}

	// user, err := userRepo.Update(&tempUser)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(user)

	// password, err := passwordRepo.Delete("OPh7lTo5t1")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(password)
	userHandler := handler.NewUserHandler(userService, redisClient)
	http.HandleFunc("/add", userHandler.HandleInitAddUser)
	http.HandleFunc("/verify", userHandler.HandleVerifyOTP)
	http.HandleFunc("/finish", userHandler.HandleFinishAddUser)
	http.ListenAndServe(":8080", nil)
}
