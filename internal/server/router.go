package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"template-go.com/internal/handlers"
)

func NewRouter(db *mongo.Database) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Health
	r.GET("/healthz", handlers.HealthHandler(db))

	// Items (example resource)
	items := r.Group("/items")
	{
		items.GET("", handlers.ListItems(db))
		items.POST("", handlers.CreateItem(db))
	}

	return r
}
