package main

import (
	"digikala/database"
	"digikala/handlers"
	"log"

	"github.com/labstack/echo/v4"
)



func main() {
	log.Println("main function called")
	e := echo.New()

	database.InitDB()
	log.Println("Connected to DataBase")

	e.POST("/Register", handlers.CreateCustomer)
	e.GET("/Customers", handlers.ReadCustomers)
	e.POST("/UpdateCustomer/:id", handlers.UpdateCustomer)

	e.Logger.Fatal(e.Start(":8080"))
}