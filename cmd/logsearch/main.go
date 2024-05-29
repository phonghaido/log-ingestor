package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phonghaido/log-ingestor/db"
	"github.com/phonghaido/log-ingestor/types"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", HandleGetHome)
	e.POST("/search", HandlePostSearch)

	e.Logger.Fatal(e.Start(":8080"))
}

func HandleGetHome(c echo.Context) error {
	return c.File("./logsearch/templates/index.html")
}

func HandlePostSearch(c echo.Context) error {
	level := c.FormValue("level")
	resourceID := c.FormValue("resourceID")
	traceID := c.FormValue("traceID")
	spanID := c.FormValue("spanID")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")

	searchData := types.SearchData{
		Level:      level,
		ResourceID: resourceID,
		TraceID:    traceID,
		SpanID:     spanID,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	mongoClient, err := db.ConnectToMongoDB(c.Request().Context())
	if err != nil {
		return err
	}

	results, err := db.Search(c, mongoClient, searchData)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("./logsearch/templates/result.html")
	if err != nil {
		return err
	}

	var renderResults bytes.Buffer
	for _, result := range results {
		err = tmpl.ExecuteTemplate(&renderResults, "result", result)
		if err != nil {
			return err
		}
	}
	return c.HTML(http.StatusOK, renderResults.String())
}
