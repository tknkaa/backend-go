package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"

	. "myapp/model"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!\n")
	})

	e.GET("/product/:id", func(c echo.Context) error {
		id := c.Param("id")
		var product Product
		result := db.First(&product, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return c.NoContent(http.StatusNotFound)
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}
		return c.JSON(http.StatusOK, product)
	})

	e.GET("/products", func(c echo.Context) error {
		var products []Product
		result := db.Find(&products)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}
		return c.JSON(http.StatusOK, products)
	})

	e.POST("product", func(c echo.Context) error {
		product := new(Product)
		if err := c.Bind(product); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		result := db.Create(product)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}
		return c.JSON(http.StatusCreated, product)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
