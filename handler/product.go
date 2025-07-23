package handler

import (
	"net/http"

	"myapp/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, world!\n")
}

func (h *Handler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	var product model.Product
	result := h.DB.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

func (h *Handler) GetAllProducts(c echo.Context) error {
	var products []model.Product
	result := h.DB.Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c echo.Context) error {
	product := new(model.Product)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	result := h.DB.Create(product)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}
	return c.JSON(http.StatusCreated, product)
}
