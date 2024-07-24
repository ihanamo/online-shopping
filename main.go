package main

import (
	"digikala/Middleware"
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

	// Customer
	e.POST("/Register", handlers.CreateCustomer)
	e.POST("/Login", handlers.LoginCustomer)

	r := e.Group("")
	r.Use(Middleware.JWTMiddleware())
	r.Use(Middleware.ExtractClaims)

	r.GET("/Customers", handlers.ReadCustomers)
	r.GET("/CustomerInfo/:id", handlers.ReadCustomer)
	r.PUT("/UpdateCustomer/:id", handlers.UpdateCustomer)
	r.DELETE("/DeleteCustomer/:id", handlers.DeleteCustomer)

	// Product
	// e.POST("/AddProduct", handlers.AddProduct)
	// e.PUT("/UpdateProduct/:id", handlers.UpdateProduct)
	// e.GET("/AllProducts", handlers.AllProducts)
	// e.GET("SpecialProduct/:id", handlers.SpecialProduct)
	// e.DELETE("/DeleteProduct/:id", handlers.DeleteProduct)

	e.Logger.Fatal(e.Start(":8080"))
}
