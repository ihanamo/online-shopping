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

	s.POST("/Product/AddProduct", handlers.AddProduct)
	s.PUT("/Product/UpdateProduct/:id", handlers.UpdateProduct)
	s.DELETE("/Product/DeleteProduct/:id", handlers.DeleteProduct)

	e.GET("/Product/AllProducts", handlers.AllProducts)
	e.GET("/Product/SpecialProduct/:type", handlers.SpecialProduct)

	// Cart
	r.POST("/Cart/Add/:product_id", handlers.AddtoCart)
	r.DELETE("/Cart/Delete/:product_id", handlers.DeleteFromCart)
	r.GET("/Cart", handlers.GetCart)
	r.POST("/Cart/Pay", handlers.PayCart)

	e.Logger.Fatal(e.Start(":8080"))
}
