package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api"
	v1 "github.com/Benyam-S/onepay/api/v1"

	urAPIHandler "github.com/Benyam-S/onepay/api/v1/http/handler"
	urHandler "github.com/Benyam-S/onepay/client/http/handler"
	"github.com/Benyam-S/onepay/client/http/session"
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

	userHandler    *urHandler.UserHandler
	userAPIHandler *urAPIHandler.UserAPIHandler
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
	apiClientRepo := repository.NewAPIClientRepository(mysqlDB)
	apiTokenRepo := repository.NewAPITokenRepository(mysqlDB)
	userService := service.NewUserService(userRepo, passwordRepo, sessionRepo, apiClientRepo, apiTokenRepo)

	userHandler = urHandler.NewUserHandler(userService, redisClient)
	userAPIHandler = urAPIHandler.NewUserAPIHandler(userService, redisClient)
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
	mysqlDB.AutoMigrate(&api.Client{})
	mysqlDB.AutoMigrate(&api.Token{})

	mysqlDB.Model(&entity.UserPassword{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&session.ServerSession{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&api.Client{}).AddForeignKey("client_user_id", "users(user_id)", "CASCADE", "CASCADE")
	mysqlDB.Model(&api.Token{}).AddForeignKey("api_key", "api_clients(api_key)", "CASCADE", "CASCADE")
}

func main() {

	configFilesDir = "C:/Users/Administrator/go/src/github.com/Benyam-S/onepay/config"

	// Initializing the server
	initServer()
	defer mysqlDB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/add", userHandler.HandleInitAddUser)
	router.HandleFunc("/verify", userHandler.HandleVerifyOTP)
	router.HandleFunc("/finish", userHandler.HandleFinishAddUser)
	router.HandleFunc("/login", userHandler.HandleLogin)
	router.HandleFunc("/filecheck", checkHandler)
	router.HandleFunc("/dashboard", tools.MiddlewareFactory(userHandler.HandleDashboard, userHandler.Authorization, userHandler.SessionDEValidation, userHandler.SessionAuthentication))
	router.HandleFunc("/logout", tools.MiddlewareFactory(userHandler.HandleLogout, userHandler.Authorization, userHandler.SessionAuthentication))

	v1.Start(userAPIHandler, router)

	http.ListenAndServe(":8080", router)
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
