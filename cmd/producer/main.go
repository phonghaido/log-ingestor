package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phonghaido/log-ingestor/cmd/producer/handlers"
	"github.com/phonghaido/log-ingestor/data"
	"github.com/phonghaido/log-ingestor/helpers"
)

var (
	kafkaConfig   = data.NewKafkaConfig("localhost:9092", "logs", 1, 1)
	kafkaProducer = handlers.NewKafkaProducer(*kafkaConfig)
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if err := kafkaProducer.CreateKafkaTopic(); err != nil {
		log.Printf("Failed to create topic %s", err)
	}

	log.Println("Server listing on port 3000...")

	e.POST("/log", helpers.ErrorWrapper(HandlePostLog))

	e.Logger.Fatal(e.Start(":3000"))
}

func HandlePostLog(c echo.Context) error {
	var reqPayload data.LogData

	if err := c.Bind(&reqPayload); err != nil {
		return helpers.InvalidJSON(c)
	}
	if err := kafkaProducer.ProduceLogKafka(reqPayload); err != nil {
		return err
	}

	return helpers.WriteJSON(c, http.StatusOK, "log ingested successfully")
}
