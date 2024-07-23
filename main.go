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

	// Customer
	e.POST("/Register", handlers.CreateCustomer)
	e.GET("/Customers", handlers.ReadCustomers)
	e.GET("/CustomerInfo/:id", handlers.ReadCustomer)
	e.PUT("/UpdateCustomer/:id", handlers.UpdateCustomer)
	e.DELETE("/DeleteCustomer/:id", handlers.DeleteCustomer)
	e.POST("Login", handlers.LoginCustomer)

	// Product
	e.POST("/AddProduct", handlers.AddProduct)
	e.PUT("/UpdateProduct/:id", handlers.UpdateProduct)
	e.GET("/AllProducts", handlers.AllProducts)
	e.GET("SpecialProduct/:id", handlers.SpecialProduct)
	e.DELETE("/DeleteProduct/:id", handlers.DeleteProduct)

	e.Logger.Fatal(e.Start(":8080"))
}
