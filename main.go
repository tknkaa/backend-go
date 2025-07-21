package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

type Product struct {
	gorm.Model
	Code  string
	Price int
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	e := echo.New()

	e.GET("/product", func(c echo.Context) error {
		var product Product
		db.First(&product, 1)
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/product", func(c echo.Context) error {
		db.Create(&Product{Code: "d42", Price: 100})
		return c.String(http.StatusOK, "added new product!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
