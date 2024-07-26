package Middleware

import (
	"digikala/DataBase"
	"digikala/models"
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

func ExtractClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok || user == nil {
			log.Println("JWT token missing or malformed")
			return c.JSON(http.StatusUnauthorized, "Missing or malformed JWT")
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid JWT claim")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		username, ok := claims["username"].(string)
		if !ok {
			log.Println("Invalid username in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		phone, ok := claims["phone"].(string)
		if !ok {
			log.Println("Invalid phone in JWT claims")
			return c.JSON(http.StatusUnauthorized, "Invalid JWT claims")
		}

		log.Printf("Extracted username: %s, Phone; %s", username, phone)

		c.Set("username", username)
		c.Set("phone", phone)

		return next(c)
	}
}


func RevokeToken(c echo.Context) error {
	tokenID := c.Param("id")
	var token models.Token
	if result := DataBase.DB.First(&token, tokenID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Token not found"})
	}

	if result := DataBase.DB.Delete(&token); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Token revoked successfully"})
}