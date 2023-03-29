package main

import (
	"clase19/cmd/server/handler"
	"clase19/internal/product"
	"clase19/pkg/middleware"
	"clase19/pkg/store"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading .env file")
	}

	storage := store.NewStore("/Users/nicolasdias/Desktop/CTD/B3/Back/C16/products.json")

	repo := product.NewRepository(storage)
	service := product.NewService(repo)
	productHandler := handler.NewProductHandler(service)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetAll())
		products.GET(":id", productHandler.GetByID())
		products.GET("/search", productHandler.Search())
		products.GET("/consumer_price", productHandler.ConsumerPrice())
		products.POST("", middleware.Authentication(), productHandler.Post())
		products.PUT(":id", middleware.Authentication(), productHandler.Put())
		products.PATCH(":id", middleware.Authentication(), productHandler.Patch())
		products.DELETE(":id", middleware.Authentication(), productHandler.Delete())
	}
	r.Run(":8080")
}
