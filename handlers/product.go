package handlers

import (
	"digikala/database"
	"digikala/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return err
	}

	result := database.DB.Create(product)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, product)
}

func AllProducts(c echo.Context) error {
	var products []models.Product
	result := database.DB.Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, products)
}

func SpecialProduct(c echo.Context) error {
	productType := c.Param("type")
	var products models.Product

	if result := database.DB.Where("type=?", productType).Find(&products); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Products not found"})
	}

	return c.JSON(http.StatusOK, products)
}

func UpdateProduct(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	updatedProduct := new(models.Product)
	if err := c.Bind(updatedProduct); err != nil {
		return err
	}

	if updatedProduct.Name != "" {
		product.Name = updatedProduct.Name
	}

	if updatedProduct.Type != "" {
		product.Type = updatedProduct.Type
	}

	if updatedProduct.Price != 0 {
		product.Price = updatedProduct.Price
	}

	if updatedProduct.Stock != 0.0 {
		product.Stock = updatedProduct.Stock
	}

	if result := database.DB.Save(&product); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	productID := c.Param("id")

	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "product not found"})
	}

	if result := database.DB.Delete(&product); result != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product deleted successfuly"})
}
