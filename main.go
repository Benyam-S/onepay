package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/onepay/api"
	v1 "github.com/Benyam-S/onepay/api/v1"
	urAPIHandler "github.com/Benyam-S/onepay/api/v1/http/handler"
	"github.com/Benyam-S/onepay/app"
	urHandler "github.com/Benyam-S/onepay/client/http/handler"
	"github.com/Benyam-S/onepay/client/http/session"
	delRepository "github.com/Benyam-S/onepay/deleted/repository"
	delService "github.com/Benyam-S/onepay/deleted/service"
	"github.com/Benyam-S/onepay/entity"
	hisRepository "github.com/Benyam-S/onepay/history/repository"
	hisService "github.com/Benyam-S/onepay/history/service"
	linkRepository "github.com/Benyam-S/onepay/linkedaccount/repository"
	linkService "github.com/Benyam-S/onepay/linkedaccount/service"
	"github.com/Benyam-S/onepay/logger"
	mtRepository "github.com/Benyam-S/onepay/moneytoken/repository"
	mtService "github.com/Benyam-S/onepay/moneytoken/service"
	"github.com/Benyam-S/onepay/tools"
	urRepository "github.com/Benyam-S/onepay/user/repository"
	urService "github.com/Benyam-S/onepay/user/service"
	walRepository "github.com/Benyam-S/onepay/wallet/repository"
	walService "github.com/Benyam-S/onepay/wallet/service"
	"github.com/go-redis/redis"

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

	onepay *app.OnePay
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
	sysConfigData, err := ioutil.ReadFile(sysConfigDir)

	// Reading data from config.onepay.json file
	onepayConfig := make(map[string]interface{})
	onepayConfigDir := filepath.Join(configFilesDir, "/config.onepay.json")
	onepayConfigData, err := ioutil.ReadFile(onepayConfigDir)

	err = json.Unmarshal(sysConfigData, &sysConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(onepayConfigData, &onepayConfig)
	if err != nil {
		panic(err)
	}

	transactionFee, ok1 := onepayConfig["transaction_fee"].(float64)
	transactionBaseLimit, ok2 := onepayConfig["transaction_base_limit"].(float64)
	withdrawBaseLimit, ok3 := onepayConfig["withdraw_base_limit"].(float64)
	dailyTransactionLimit, ok4 := onepayConfig["daily_send_limit"].(float64)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		panic(errors.New("unable to parse onepay config data"))
	}

	// Setting enviromental variables so they can be used any where on the application
	os.Setenv("config_files_dir", configFilesDir)
	os.Setenv("onepay_secret_key", sysConfig.SecretKey)
	os.Setenv("onepay_cookie_name", sysConfig.CookieName)

	os.Setenv(entity.TransactionFee, fmt.Sprintf("%f", transactionFee))
	os.Setenv(entity.TransactionBaseLimit, fmt.Sprintf("%f", transactionBaseLimit))
	os.Setenv(entity.WithdrawBaseLimit, fmt.Sprintf("%f", withdrawBaseLimit))
	os.Setenv(entity.DailyTransactionLimit, fmt.Sprintf("%f", dailyTransactionLimit))

	// Initializing the database with the needed tables and values
	initDB()

	userRepo := urRepository.NewUserRepository(mysqlDB)
	passwordRepo := urRepository.NewPasswordRepository(mysqlDB)
	sessionRepo := urRepository.NewSessionRepository(mysqlDB)
	apiClientRepo := urRepository.NewAPIClientRepository(mysqlDB)
	apiTokenRepo := urRepository.NewAPITokenRepository(mysqlDB)
	walletRepo := walRepository.NewWalletRepository(mysqlDB)
	historyRepo := hisRepository.NewHistoryRepository(mysqlDB)
	linkedAccountRepo := linkRepository.NewLinkedAccountRepository(mysqlDB)
	moneyTokenRepo := mtRepository.NewMoneyTokenRepository(mysqlDB)
	deletedUserRepo := delRepository.NewDeletedUserRepository(mysqlDB)
	deletedLinkedAccountRepo := delRepository.NewDeletedLinkedAccountRepository(mysqlDB)

	userService := urService.NewUserService(userRepo, passwordRepo, sessionRepo, apiClientRepo, apiTokenRepo)
	deletedService := delService.NewDeletedService(deletedUserRepo, deletedLinkedAccountRepo)
	walletService := walService.NewWalletService(walletRepo)
	historyService := hisService.NewHistoryService(historyRepo)
	linkedAccountService := linkService.NewLinkedAccountService(linkedAccountRepo)
	moneyTokenService := mtService.NewMoneyTokenService(moneyTokenRepo)

	path, _ := os.Getwd()
	path = filepath.Join(path, "./logger")
	dataLogger := logger.NewLogger(path)
	channel := make(chan string)

	onepay = app.NewApp(walletService, historyService, linkedAccountService, moneyTokenService,
		dataLogger, channel)

	userHandler = urHandler.NewUserHandler(userService, redisClient)
	userAPIHandler = urAPIHandler.NewUserAPIHandler(onepay, userService, deletedService, redisClient)
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
	mysqlDB.AutoMigrate(&entity.UserHistory{})
	mysqlDB.AutoMigrate(&entity.UserWallet{})
	mysqlDB.AutoMigrate(&entity.MoneyToken{})
	mysqlDB.AutoMigrate(&entity.LinkedAccount{})
	mysqlDB.AutoMigrate(&entity.DeletedUser{})
	mysqlDB.AutoMigrate(&entity.DeletedLinkedAccount{})

	// ----------------------------- This must be changed -----------------------------
	count := 0
	mysqlDB.Exec("CREATE TABLE IF NOT EXISTS extras (total_user_count int)")
	mysqlDB.Table("extras").Count(&count)
	if count != 1 {
		mysqlDB.Exec("DELETE FROM extras")
		mysqlDB.Exec("INSERT INTO extras VALUES (0)")
	}
	// --------------------------------------------------------------------------------
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
	router.HandleFunc("/dashboard", tools.MiddlewareFactory(userHandler.HandleDashboard, userHandler.Authorization, userHandler.SessionDEValidation, userHandler.SessionAuthentication))
	router.HandleFunc("/logout", tools.MiddlewareFactory(userHandler.HandleLogout, userHandler.Authorization, userHandler.SessionAuthentication))

	v1.Start(userAPIHandler, router)

	go func() {
		for {
			time.Sleep(time.Minute * 30)
			onepay.Channel <- "all"
		}
	}()

	go func() {

		for {

			value := <-onepay.Channel
			switch value {

			case "all":
				onepay.ReloadMoneyToken()
				onepay.ReloadWallet()
				onepay.ReloadHistory()

			case "reload_money_token":
				onepay.ReloadMoneyToken()

			case "reload_wallet":
				onepay.ReloadWallet()
				fallthrough

			case "reload_history":
				onepay.ReloadHistory()
			}
		}
	}()

	http.ListenAndServe(":8080", router)
}
