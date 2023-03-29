package main

import (
	"clase18/cmd/server/handler"
	"clase18/internal/product"
	"clase18/pkg/store"
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetAll())
		products.GET(":id", productHandler.GetByID())
		products.GET("/search", productHandler.Search())
		products.GET("/consumer_price", productHandler.ConsumerPrice())
		products.POST("", productHandler.Post())
		products.PUT(":id", productHandler.Put())
		products.PATCH(":id", productHandler.Patch())
		products.DELETE(":id", productHandler.Delete())
	}
	r.Run(":8080")
}
