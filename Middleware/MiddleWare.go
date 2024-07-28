package Middleware

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

var JWTKey = []byte("secret_key")

func JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey: JWTKey,
		ContextKey: "user",
	}

	return echojwt.WithConfig(config)
}

func ExtractClaimsCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok || user == nil {
			log.Println("JWT token missing or malformed")
			return c.JSON(http.StatusUnauthorized, "Missing or malformed JWT")
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid JWT claims structure")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims structure")
		}

		customerID, ok := claims["customer-id"].(float64)
		if !ok {
			log.Println("Invalid customer-id in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		username, ok := claims["username"].(string)
		if !ok {
			log.Println("Invalid username in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		log.Printf("Extracted customer-id: %f, username: %s", customerID, username)

		c.Set("customer-id", uint(customerID))
		c.Set("username", username)

		return next(c)
	}
}

func ExtractClaimsManager(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok || user == nil {
			log.Println("JWT token missing or malformed")
			return c.JSON(http.StatusUnauthorized, "Missing or malformed JWT")
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid JWT claims structure")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims structure")
		}

		managerID, ok := claims["manager-id"].(float64)
		if !ok {
			log.Println("Invalid manager-id in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		username, ok := claims["username"].(string)
		if !ok {
			log.Println("Invalid username in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		log.Printf("Extracted manager-id: %f, username: %s", managerID, username)

		c.Set("manager-id", uint(managerID))
		c.Set("username", username)

		return next(c)
	}
}
