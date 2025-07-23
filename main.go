package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"myapp/handler"
	"myapp/model"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Product{}, &model.User{}, &model.Session{})

	e := echo.New()

	h := &handler.Handler{DB: db}

	e.GET("/", h.GetRoot)
	e.GET("/product/:id", h.GetProductByID)
	e.GET("/products", h.GetAllProducts)
	e.POST("/product", h.CreateProduct)

	e.POST("/sign-up", h.SignUp)

	e.Logger.Fatal(e.Start(":1323"))
}
