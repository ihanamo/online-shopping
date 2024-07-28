package handlers

import (
	"digikala/DataBase"
	"digikala/models"
	"log"
	"net/http"
	"strconv"
	"time"

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
		log.Println("Invalid JWT claims structure")
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	customerID, ok := claims["customer-id"].(float64)
	if !ok {
		log.Println("Invalid customer-id in JWT claims")
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	productID := c.Param("id")
	var product models.Product
	if result := DataBase.DB.First(&product, productID); result.Error != nil {
		log.Println(productID)
		log.Println("Error finding product:", result.Error)
		log.Println("Product found", product)
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

	var cartProduct models.CartProduct
	if result := DataBase.DB.Where("cart_id = ? AND product_id = ?", cart.ID, product.ID).First(&cartProduct); result.Error != nil {
		cartProduct = models.CartProduct{
			CartID:    cart.ID,
			ProductID: product.ID,
			Quantity:  1,
		}
		if result := DataBase.DB.Create(&cartProduct); result.Error != nil {
			log.Println("Error creating CartProduct:", result.Error)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to add product to cart"})
		}
	} else {
		cartProduct.Quantity += 1
		if result := DataBase.DB.Save(&cartProduct); result.Error != nil {
			log.Println("Error updating CartProduct:", result.Error)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update cart product"})
		}
	}

	cart.Total += product.Price
	if result := DataBase.DB.Save(&cart); result.Error != nil {
		log.Println("Error updating Cart:", result.Error)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update cart"})
	}

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
		log.Println("Invalid JWT claims structure")
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	customerID, ok := claims["customer-id"].(float64)
	if !ok {
		log.Println("Invalid customer-id in JWT claims")
		return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
	}

	productIDParam := c.Param("id")
	productID, err := strconv.ParseUint(productIDParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid product ID")
	}

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
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Active cart not found"})
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

		payment := models.Payment{
			PaidAt:     time.Now(),
			Total:      cart.Total,
			CustomerID: cart.CustomerID,
		}
		if result := DataBase.DB.Create(&payment); result.Error != nil {
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

func DeleteCart(c echo.Context) error {
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
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Active cart not found"})
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

	if result := DataBase.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartProduct{}); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	if result := DataBase.DB.Delete(&cart); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Cart deleted successfully"})
}
