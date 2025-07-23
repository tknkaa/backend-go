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

	e.Static("/", "public")

	e.POST("/sign-up", h.SignUp)
	e.POST("/sign-in", h.SignIn)

	api := e.Group("/api")
	api.Use(h.AuthMiddleware)
	api.GET("/", h.GetRoot)
	api.GET("/product/:id", h.GetProductByID)
	api.GET("/products", h.GetAllProducts)
	api.POST("/product", h.CreateProduct)

	e.Logger.Fatal(e.Start(":1323"))
}
