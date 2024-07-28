package handlers

import (
	"digikala/DataBase"
	"digikala/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var JWTKey = []byte("secret_key")

func GenerateJWTCustomer(customer models.Customer) (string, error) {
	claims := &jwt.MapClaims{
		"customer-id": customer.ID,
		"username":    customer.Username,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func AuthenticateCustomer(username, password string) (models.Customer, string, error) {
	var customer models.Customer
	log.Println("Authenticating user:", username)
	result := DataBase.DB.Where("username = ?", username).First(&customer)
	if result.Error != nil {
		log.Println("User not found:", result.Error)
		return customer, "", echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return customer, "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
	}

	token, err := GenerateJWTCustomer(customer)
	if err != nil {
		return customer, "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to Generate Token")
	}

	return customer, token, nil
}

func CreateCustomer(c echo.Context) error {
	log.Println("Create Customer called")
	customer := new(models.Customer)
	if err := c.Bind(customer); err != nil {
		return err
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"messgae": "Failed to hash password"})
	}
	customer.Password = string(hashPass)
	log.Println("the hash password is:", customer.Password)

	result := DataBase.DB.Create(customer)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	log.Println("user created")

	token, err := GenerateJWTCustomer(*customer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfuly",
		"user":    customer,
		"token":   token,
	})
}

func LoginCustomer(c echo.Context) error {
	credentials := new(models.Credentials)
	if err := c.Bind(credentials); err != nil {
		log.Println("Failed to bind credentials:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request data"})
	}

	log.Println("Credentials received:", credentials)
	customer, token, err := AuthenticateCustomer(credentials.Username, credentials.Password)
	if err != nil {
		log.Println("Authentication failed:", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   token,
		"user":    customer,
	})
}

func ReadCustomer(c echo.Context) error {
	customerID := c.Param("id")
	var customer models.Customer
	if result := DataBase.DB.First(&customer, customerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	return c.JSON(http.StatusOK, customer)
}

func ReadCustomers(c echo.Context) error {
	var customers []models.Customer
	result := DataBase.DB.Find(&customers)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, customers)
}

func UpdateCustomer(c echo.Context) error {
	customerID := c.Param("id")
	var customer models.Customer
	if result := DataBase.DB.First(&customer, customerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	updatedCustomer := new(models.Customer)
	if err := c.Bind(updatedCustomer); err != nil {
		return err
	}

	if updatedCustomer.FirstName != "" {
		customer.FirstName = updatedCustomer.FirstName
	}

	if updatedCustomer.LastName != "" {
		customer.LastName = updatedCustomer.LastName
	}

	if updatedCustomer.Phone != "" {
		customer.Phone = updatedCustomer.Phone
	}

	if updatedCustomer.Password != "" {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(updatedCustomer.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		customer.Password = string(hashPass)
	}

	if result := DataBase.DB.Save(&customer); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c echo.Context) error {
	customerID := c.Param("id")

	var customer models.Customer
	if result := DataBase.DB.First(&customer, customerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "user not found"})
	}

	if result := DataBase.DB.Delete(&customer); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfuly"})

}
