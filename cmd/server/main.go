package main

import (
	"clase19/cmd/server/handler"
	"clase19/docs"
	"clase19/internal/product"
	"clase19/pkg/middleware"
	"clase19/pkg/store"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Products Market
// @version 1.0
// @description This API Handle Products.
// @termsOfService https://developers.ctd.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.ctd.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		panic(err.Error())
	}
	errPing := db.Ping()
	if errPing != nil {
		panic(errPing.Error())
	}
	storage := store.NewSqlStore(db)
	// storage := store.NewJsonStore("../../products.json")

	repo := product.NewRepository(storage)
	service := product.NewService(repo)
	productHandler := handler.NewProductHandler(service)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

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
