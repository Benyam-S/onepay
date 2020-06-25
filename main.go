package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Benyam-S/onepay/session"

	"github.com/Benyam-S/onepay/client/http/handler"
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user/service"
	"github.com/go-redis/redis"

	"github.com/Benyam-S/onepay/user/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	configFilesDir string
	redisClient    *redis.Client
	mysqlDB        *gorm.DB
	sysConfig      SystemConfig
	err            error

	userHandler *handler.UserHandler
)

// SystemConfig is a type that defines a server system configuration file
type SystemConfig struct {
	RedisClient map[string]string `json:"redis_client"`
	MysqlClient map[string]string `json:"mysql_client"`
	CookieName  string            `json:"cookie_name"`
	SecretKey   string            `json:"secret_key"`
}

// initServer initialize the web server for takeoff
func initServer() {

	// Reading data from config.server.json file and creating the systemconfig  object
	sysConfigDir := filepath.Join(configFilesDir, "/config.server.json")
	data, err := ioutil.ReadFile(sysConfigDir)

	err = json.Unmarshal(data, &sysConfig)
	if err != nil {
		panic(err)
	}

	// Setting enviromental variables so they can be used any where on the application
	os.Setenv("config_files_dir", configFilesDir)
	os.Setenv("onepay_secret_key", sysConfig.SecretKey)
	os.Setenv("onepay_cookie_name", sysConfig.CookieName)

	// Initializing the database with the needed tables and values
	initDB()

	userRepo := repository.NewUserRepository(mysqlDB)
	passwordRepo := repository.NewPasswordRepository(mysqlDB)
	sessionRepo := repository.NewSessionRepository(mysqlDB)
	userService := service.NewUserService(userRepo, passwordRepo, sessionRepo)
	userHandler = handler.NewUserHandler(userService, redisClient)
}

// initDB initialize the database for takeoff
func initDB() {

	redisDB, _ := strconv.ParseInt(sysConfig.RedisClient["database"], 0, 0)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     sysConfig.RedisClient["address"] + ":" + sysConfig.RedisClient["port"],
		Password: sysConfig.RedisClient["password"], // no password set
		DB:       int(redisDB),                      // use default DB
	})

	mysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sysConfig.MysqlClient["user"], sysConfig.MysqlClient["password"],
		sysConfig.MysqlClient["address"], sysConfig.MysqlClient["port"], sysConfig.MysqlClient["database"]))

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database: mysql @GORM")

	// Creating and Migrating tables from the structurs
	mysqlDB.AutoMigrate(&entity.UserPassword{})
	mysqlDB.AutoMigrate(&entity.User{})
	mysqlDB.AutoMigrate(&session.ServerSession{})
	mysqlDB.Model(&entity.UserPassword{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
}

func main() {

	configFilesDir = "C:/Users/Administrator/go/src/github.com/Benyam-S/onepay/config"

	// Initializing the server
	initServer()
	defer mysqlDB.Close()

	http.HandleFunc("/add", userHandler.HandleInitAddUser)
	http.HandleFunc("/verify", userHandler.HandleVerifyOTP)
	http.HandleFunc("/finish", userHandler.HandleFinishAddUser)
	http.HandleFunc("/login", userHandler.HandleLogin)
	http.HandleFunc("/filecheck", checkHandler)
	http.HandleFunc("/dashboard", tools.MiddlewareFactory(userHandler.HandleDashboard, userHandler.Authorization, userHandler.SessionDEValidation, userHandler.SessionAuthentication))
	http.HandleFunc("/logout", tools.MiddlewareFactory(userHandler.HandleLogout, userHandler.Authorization, userHandler.SessionAuthentication))
	http.ListenAndServe(":8080", nil)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./hello.html")
}

func checkFunc() {
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
}
