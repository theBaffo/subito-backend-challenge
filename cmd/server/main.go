package main

import (
	"log"
	"net/http"

	"github.com/theBaffo/subito-backend-challenge/internal/handler"
	"github.com/theBaffo/subito-backend-challenge/internal/middleware"
	"github.com/theBaffo/subito-backend-challenge/internal/repository/memory"
	"github.com/theBaffo/subito-backend-challenge/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// --- Dependency wiring ---
	// To swap storage: replace memory stores with a PostgreSQL implementation
	// that satisfies the same repository interfaces.
	productStore := memory.NewProductStore()
	orderStore := memory.NewOrderStore()

	productSvc := service.NewProductService(productStore)
	orderSvc := service.NewOrderService(orderStore, productStore)

	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// --- Router setup ---
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/products", productHandler.ListProducts)
		v1.GET("/products/:id", productHandler.GetProduct)

		v1.POST("/orders", orderHandler.CreateOrder)
		v1.GET("/orders/:id", orderHandler.GetOrder)
	}

	log.Println("Purchase Cart Service starting on :8080")
	
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
