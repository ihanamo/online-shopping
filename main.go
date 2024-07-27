package main

import (
	"digikala/DataBase"
	"digikala/Middleware"
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
	e.POST("/RegisterCustomer", handlers.CreateCustomer)
	e.POST("/LoginCustomer", handlers.LoginCustomer)

	r := e.Group("")
	r.Use(Middleware.JWTMiddleware())
	r.Use(Middleware.ExtractClaims)

	r.GET("/Customers", handlers.ReadCustomers)
	r.GET("/CustomerInfo/:id", handlers.ReadCustomer)
	r.PUT("/UpdateCustomer/:id", handlers.UpdateCustomer)
	r.DELETE("/DeleteCustomer/:id", handlers.DeleteCustomer)

	// Product
	e.POST("/RegisterManager", handlers.CreateManager)
	e.POST("/LoginManager", handlers.LoginManager)

	s := e.Group("")
	s.Use(Middleware.JWTMiddleware())
	s.Use(Middleware.ExtractClaims)

	s.POST("/AddProduct", handlers.AddProduct)
	s.PUT("/UpdateProduct/:id", handlers.UpdateProduct)
	s.DELETE("/DeleteProduct/:id", handlers.DeleteProduct)

	e.GET("/AllProducts", handlers.AllProducts)
	e.GET("SpecialProduct/:type", handlers.SpecialProduct)

	// Cart
	r.POST("/Cart/Add/:product_id", handlers.AddtoCart)

	e.Logger.Fatal(e.Start(":8080"))
}
