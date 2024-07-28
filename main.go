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
	e.POST("/Customer/RegisterCustomer", handlers.CreateCustomer)
	e.POST("/Customer/LoginCustomer", handlers.LoginCustomer)

	r := e.Group("")
	r.Use(Middleware.JWTMiddleware())
	r.Use(Middleware.ExtractClaimsCustomer)

	r.GET("/Customer/AllCustomers", handlers.ReadCustomers)
	r.GET("/Customer/CustomerInfo/:id", handlers.ReadCustomer)
	r.PUT("/Customer/UpdateCustomer/:id", handlers.UpdateCustomer)
	r.DELETE("/Customer/DeleteCustomer/:id", handlers.DeleteCustomer)

	// Product
	e.POST("/Manager/RegisterManager", handlers.CreateManager)
	e.POST("/Manager/LoginManager", handlers.LoginManager)

	s := e.Group("")
	s.Use(Middleware.JWTMiddleware())
	s.Use(Middleware.ExtractClaimsManager)

	s.POST("/Product/AddProduct", handlers.AddProduct)
	s.PUT("/Product/UpdateProduct/:id", handlers.UpdateProduct)
	s.DELETE("/Product/DeleteProduct/:id", handlers.DeleteProduct)

	e.GET("/Product/AllProducts", handlers.AllProducts)
	e.GET("/Product/SpecialProduct/:type", handlers.SpecialProduct)

	// Cart
	r.POST("/Cart/Add/:id", handlers.AddtoCart)
	r.DELETE("/Cart/DeleteProduct/:id", handlers.DeleteFromCart)
	r.GET("/Cart", handlers.GetCart)
	r.POST("/Cart/Pay", handlers.PayCart)
	r.DELETE("/Cart/DeleteCart", handlers.DeleteCart)

	// Payment
	r.GET("/Payment", handlers.GetAllPayments)

	e.Logger.Fatal(e.Start(":8080"))

}
