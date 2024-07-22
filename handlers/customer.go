package handlers

import (
	"digikala/database"
	"digikala/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

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

	result := database.DB.Create(customer)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	log.Println("user created")
	return c.JSON(http.StatusCreated, customer)
}

// func Login(c echo.Context) error {
	
// }

func ReadCustomers(c echo.Context) error {
	var customers []models.Customer
	result := database.DB.Find(&customers)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, customers)
}

func UpdateCustomer(c echo.Context) error {
	customerID := c.Param("id")
	var customer models.Customer
	if result := database.DB.First(&customer, customerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	updatedCustomer := new(models.Customer)
	if err := c.Bind(updatedCustomer); err != nil {
		return err
	}

	if updatedCustomer.FirsName != "" {
		customer.FirsName = updatedCustomer.FirsName
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

	if result := database.DB.Save(&customer); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c echo.Context) error {
	customerID := c.Param("id")

	var customer models.Customer
	if result := database.DB.First(&customer, customerID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "user not found"})
	}

	if result := database.DB.Delete(&customer); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfuly"})

}
