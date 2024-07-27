package handlers

import (
	"digikala/DataBase"
	"digikala/models"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AddtoCart(c echo.Context) error {
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

	productID := c.Param("product_id")
	var product models.Product
	if result := DataBase.DB.First(&product, productID); result != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Product not found"})
	}

	if product.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Product out of stock"})
	}

	var cart models.Cart
	if result := DataBase.DB.Where("customer_id = ? AND is_payed = ?", uint(customerID), false).First(&cart); result.Error != nil {
		cart = models.Cart{
			CustomerID: uint(customerID),
			Total:      0,
			IsPayed:    false,
		}
		DataBase.DB.Create(&cart)
	}

	cartProduct := models.CartProduct{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  1,
	}
	DataBase.DB.Create(&cartProduct)

	cart.Total += product.Price
	DataBase.DB.Save(&cart)

	product.Stock -= 1
	if result := DataBase.DB.Save(&product); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update product stock"})
	}

	return c.JSON(http.StatusOK, cart)
}

func DeleteFromCart(c echo.Context) error {
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

	productID := c.Param("product_id")
	var product models.Product
	if result := DataBase.DB.First(&product, productID); result != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Product not found"})
	}

	var cart models.Cart
	if result := DataBase.DB.Where("customer_id = ? AND is_payed = ?", uint(customerID)).First(&cart); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Cart not found"})
	}

	var cartProduct models.CartProduct
	if result := DataBase.DB.Where("cart_id = ? ANd product_id = ?", cart.ID, product.ID).First(&cartProduct); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Product not found in cart"})
	}

	if result := DataBase.DB.Delete(&cartProduct); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	cart.Total -= product.Price * float64(cartProduct.Quantity)
	if result := DataBase.DB.Save(&cart); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	product.Stock += cartProduct.Quantity
	if result := DataBase.DB.Save(&product); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product removed from cart successfuly"})
}

func GetCart(c echo.Context) error {
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

	var cart models.Cart
	if result := DataBase.DB.Where("customer_id = ? AND is_payed = ?", uint(customerID), false).Preload("Products").First(&cart); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Cart not found"})
	}

	return c.JSON(http.StatusOK, cart)
}

func PayCart(c echo.Context) error {
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

	var cart models.Cart
	if result := DataBase.DB.Where("customer_id = ? AND is_payed = ?", uint(customerID), false).Preload("Products").First(&cart); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Cart not found"})
	}

	if len(cart.Products) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Cart is empty"})
	}

	paymentSuccessful := true

	if paymentSuccessful {
		cart.IsPayed = true
		if result := DataBase.DB.Save(&cart); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		for _, product := range cart.Products {
			var cartProduct models.CartProduct
			if result := DataBase.DB.Where("cart_id = ? AND product_id = ?", cart.ID, product.ID).First(&cartProduct); result.Error != nil {
				return c.JSON(http.StatusInternalServerError, result.Error)
			}

			product.Stock -= cartProduct.Quantity
			if result := DataBase.DB.Save(&product); result.Error != nil {
				return c.JSON(http.StatusInternalServerError, result.Error)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Payment successful, cart has been paid"})
	} else {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Payment failed"})
	}
}
