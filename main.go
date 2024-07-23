package main

import (
	"digikala/DataBase"
	"digikala/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("main function called")
	e := echo.New()

	DataBase.InitDB()
	log.Println("Connected to DataBase")

	e.POST("/Register", handlers.CreateCustomer)
	e.GET("/Customers", handlers.ReadCustomers)
	e.GET("/Info/:id", handlers.ReadCustomer)
	e.PUT("/UpdateCustomer/:id", handlers.UpdateCustomer)
	e.DELETE("/DeleteCustomer/:id", handlers.DeleteCustomer)
	e.POST("Login", handlers.LoginCustomer)

	e.Logger.Fatal(e.Start(":8080"))
}
