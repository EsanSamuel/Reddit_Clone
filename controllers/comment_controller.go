package controllers

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var comment models.Comment

		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error binding subreddit payload", "details": err.Error()})
			return
		}

		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		comment.CommentID = bson.NewObjectID().Hex()

		result, err := database.CommentCollection.InsertOne(ctx, comment)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating comment", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			database.PostCollection.UpdateOne(ctx, bson.M{"post_id": comment.PostID}, bson.M{"$inc": bson.M{"comment_count": 1}})
			if comment.ParentID != "" {
				database.CommentCollection.UpdateOne(ctx, bson.M{"parent_id": comment.ParentID}, bson.M{"$inc": bson.M{"comment_count": 1}})
			}
		}

		c.JSON(http.StatusCreated, gin.H{"comment": comment, "result": result})

	}
}

func GetPostComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		post_id := c.Param("post_id")

		var comments []models.Comment

		filter := bson.M{"post_id": post_id}
		findOptions := options.Find()

		// Search Subreddit
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)
			filter["$and"] = []bson.M{
				{
					"post_id": post_id,
				},
				{
					"$or": []bson.M{
						{
							"type": bson.M{
								"$regex":   safe,
								"$options": "i",
							},
						},
						{
							"content": bson.M{
								"$regex":   safe,
								"$options": "i",
							},
						},
					},
				},
			}

		}

		// Sort Subreddit
		if sort := strings.TrimSpace(c.Query("sort")); sort != "" {
			switch sort {
			case "asc":
				findOptions.SetSort(bson.D{{"created_at", 1}})
			case "desc":
				findOptions.SetSort(bson.D{{"created_at", -1}})
			}

		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 9

		findOptions.SetSkip((int64(page) - 1) * int64(perPage))
		findOptions.SetLimit(int64(perPage))

		cursor, err := database.CommentCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting subreddit post", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &comments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding subreddit post", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		c.JSON(http.StatusCreated, comments)
	}
}

func GetParentComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		parent_id := c.Param("parent_id")

		var comments []models.Comment

		filter := bson.M{"parent_id": parent_id}
		findOptions := options.Find()

		// Search Subreddit
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)
			filter["$and"] = []bson.M{
				{
					"parent_id": parent_id,
				},
				{
					"$or": []bson.M{
						{
							"type": bson.M{
								"$regex":   safe,
								"$options": "i",
							},
						},
						{
							"content": bson.M{
								"$regex":   safe,
								"$options": "i",
							},
						},
					},
				},
			}

		}

		// Sort Subreddit
		if sort := strings.TrimSpace(c.Query("sort")); sort != "" {
			switch sort {
			case "asc":
				findOptions.SetSort(bson.D{{"created_at", 1}})
			case "desc":
				findOptions.SetSort(bson.D{{"created_at", -1}})
			}

		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 9

		findOptions.SetSkip((int64(page) - 1) * int64(perPage))
		findOptions.SetLimit(int64(perPage))

		cursor, err := database.CommentCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting subreddit post", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &comments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding subreddit post", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		c.JSON(http.StatusCreated, comments)
	}
}

func GetCommentById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		commentId := c.Param("id")

		var comment models.Comment

		err := database.CommentCollection.FindOne(ctx, bson.M{"comment_id": commentId}).Decode(&comment)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding comment", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}
