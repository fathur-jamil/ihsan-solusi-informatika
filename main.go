package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"account_service/controller"
	"account_service/model"
	"account_service/repository"
	"account_service/service"
	"account_service/util"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error load config, %v", err)
	}

	db, err := initDB(config)
	if err != nil {
		log.Fatalf("Error init db, %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	accountRepository := repository.NewAccountRepository(db)

	registrationService := service.NewRegistrationService(accountRepository, userRepository)
	transactionService := service.NewTransactionService(accountRepository, userRepository)

	registrationController := controller.NewRegistrationController(registrationService)
	transactionController := controller.NewTransactionController(transactionService)

	e := echo.New()
	e.Validator = &util.CustomValidator{Validator: validator.New()}

	e.POST("/daftar", registrationController.Register)
	e.POST("/tabung", transactionController.Deposit)
	e.POST("/tarik", transactionController.Withdraw)

	e.GET("/saldo/:no_rekening", transactionController.GetBalance)

	address := fmt.Sprintf("%s:%d", config.APIHost, config.APIPort)
	log.Infof("Server starting on: %s", address)
	e.Logger.Fatal(e.Start(address))
}

type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	APIHost    string
	APIPort    int
}

func loadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	apiHost := flag.String("host", "localhost", "REST API host")
	apiPort := flag.Int("port", 8080, "REST API port")
	flag.Parse()

	dbHost := "localhost"
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}

	return Config{
		DBDriver:   "postgres",
		DBHost:     dbHost,
		DBPort:     "5432",
		DBName:     "account",
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		APIHost:    *apiHost,
		APIPort:    *apiPort,
	}, nil
}

func initDB(config Config) (*gorm.DB, error) {
	DBURL := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable ",
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBUser,
		config.DBPassword,
	)
	db, err := gorm.Open(config.DBDriver, DBURL)
	if err != nil {
		return nil, err
	}

	//db.LogMode(true)
	db.Debug().AutoMigrate(
		&model.Account{},
		&model.User{},
	)
	return db, nil
}
