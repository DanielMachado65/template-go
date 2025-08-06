package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"template-go.com/internal/models"
)

const collectionName = "items"

type createItemDTO struct {
	Name string `json:"name" binding:"required"`
}

func ListItems(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		cur, err := db.Collection(collectionName).Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(ctx)

		var items []models.Item
		for cur.Next(ctx) {
			var it models.Item
			if err := cur.Decode(&it); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			items = append(items, it)
		}
		if err := cur.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}

func CreateItem(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto createItemDTO
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doc := models.Item{
			Name:      dto.Name,
			CreatedAt: time.Now().UTC(),
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		res, err := db.Collection(collectionName).InsertOne(ctx, doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": res.InsertedID})
	}
}
