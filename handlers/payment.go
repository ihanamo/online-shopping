package handlers

import (
	"digikala/DataBase"
	"digikala/models"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetAllPayments(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.JSON(http.StatusUnauthorized, "Missing or malformed JWT")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	customerID, ok := claims["customer-id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	var payments []models.Payment
	if result := DataBase.DB.Where("customer_id = ?", uint(customerID)).Find(&payments); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, payments)
}
