package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/adapter"
	"github.com/SaidHernandez/bia-comsumtion/business/model"
	"github.com/SaidHernandez/bia-comsumtion/business/repository"
	handlers "github.com/SaidHernandez/bia-comsumtion/handler"
	"github.com/SaidHernandez/bia-comsumtion/infraestructure/cache"
	"github.com/SaidHernandez/bia-comsumtion/infraestructure/db"
	"github.com/SaidHernandez/bia-comsumtion/services"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var consumptionHandler *handlers.ConsumptionHandler

func populateConsumptionDBFromCSV(db *gorm.DB, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	var consumptions []model.Consumption
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		consumptionID := record[0]
		if consumptionID == "" {
			return fmt.Errorf("Consumption Id Invalido: %w", err)
		}

		meterID, err := strconv.Atoi(record[1])
		if err != nil {
			return fmt.Errorf("Invalido meterID formato: %w", err)
		}
		active, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return fmt.Errorf("Invalido active energy formato: %w", err)
		}
		date, err := time.Parse("2006-01-02 15:04:05-07", record[3])
		if err != nil {
			return fmt.Errorf("Invalido date format: %w", err)
		}

		reactiveInductive := rand.Float64() * 20000
		reactiveCapacitive := rand.Float64() * 20000
		exported := rand.Float64() * 20000

		consumptions = append(consumptions, model.Consumption{
			ID:                 consumptionID,
			MeterID:            meterID,
			Date:               date,
			ActiveEnergy:       active,
			ReactiveInductive:  reactiveInductive,
			ReactiveCapacitive: reactiveCapacitive,
			ExportedEnergy:     exported,
		})

	}

	if len(consumptions) == 0 {
		return errors.New("no valid consumption records found")
	}

	batchSize := 500
	for i := 0; i < len(consumptions); i += batchSize {
		end := i + batchSize
		if end > len(consumptions) {
			end = len(consumptions)
		}

		if err := db.Create(consumptions[i:end]).Error; err != nil {
			return fmt.Errorf("failed to insert records into database: %w", err)
		}
	}

	return nil
}

func initDB() error {
	var command string

	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	var err error
	err = db.InitDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.DB.AutoMigrate(&model.Consumption{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if command == "runMigration" {
		err = populateConsumptionDBFromCSV(db.DB, "./infraestructure/resources/test_bia11.csv")
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func initServices() {

	cacheInstance := cache.NewMemoryCache()
	adapterInstance := adapter.NewAddressAdapter()
	repository := repository.NewConsumptionRepository()

	addressService := services.NewAddressServiceClient(cacheInstance, adapterInstance)
	consumptionService := services.NewConsumptionService(addressService, repository)
	consumptionHandler = handlers.NewConsumptionHandler(consumptionService)
}

func main() {
	initDB()
	initServices()

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/consumption", consumptionHandler.GetConsumption)
	e.Logger.Fatal(e.Start(":8080"))
}
