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
	"github.com/EsanSamuel/Reddit_Clone/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var post models.Post

		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		if isMultipart {
			if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid multipart form",
					"details": err.Error(),
				})
				return
			}
		}

		if err := c.ShouldBind(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "error binding post payload",
				"details": err.Error(),
			})
			return
		}

		if isMultipart {
			if form, _ := c.MultipartForm(); form != nil {
				if files, ok := form.File["files"]; ok && len(files) > 0 {
					post.FileUrls = utils.UploadFiles(c, files)
				}
			}
		}

		post.PostID = bson.NewObjectID().Hex()
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()
		post.Score = 0

		result, err := database.PostCollection.InsertOne(ctx, post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "error creating post",
				"details": err.Error(),
			})
			return
		}

		if result.Acknowledged {
			_, _ = database.SubredditCollection.UpdateOne(
				ctx,
				bson.M{"subreddit_id": post.SubredditID},
				bson.M{"$inc": bson.M{"posts_count": 1}},
			)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "post created successfully",
			"post_id": post.PostID,
		})
	}
}

func GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var posts []models.Post

		filter := bson.M{}
		findOptions := options.Find()

		// Search Subreddit
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)
			filter = bson.M{
				"$or": []bson.M{
					{
						"title": bson.M{
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

		cursor, err := database.PostCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting subreddit post", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &posts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding subreddit post", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, posts)

	}
}

func GetSubRedditPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		subreddit_id := c.Param("subreddit_id")

		filter := bson.M{"subreddit_id": subreddit_id}
		findOptions := options.Find()

		// Search posts
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)

			filter["$and"] = []bson.M{
				{
					"subreddit_id": subreddit_id,
				},
				{
					"$or": []bson.M{
						{
							"title": bson.M{
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

		var posts []models.Post

		cursor, err := database.PostCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting subreddit post", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &posts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding subreddit post", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		c.JSON(http.StatusCreated, posts)

	}
}

func GetTagPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var posts []models.Post

		cursor, err := database.PostCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting subreddit post", "details": err.Error()})
			return
		}

		if err := cursor.All(ctx, &posts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding subreddit post", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		tagMaps := make(map[string][]models.Post)
		response := make(map[string]models.TagsPosts)

		for _, post := range posts {
			for _, posttag := range post.Tags {
				tagMaps[posttag] = append(tagMaps[posttag], post)
				post_counts := len(tagMaps[posttag])
				response[posttag] = models.TagsPosts{
					Posts:     tagMaps[posttag],
					PostCount: post_counts,
				}
			}
		}
		c.JSON(http.StatusOK, response)

	}
}

func GetPostById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		postId := c.Param("id")

		var post models.Post

		err := database.PostCollection.FindOne(ctx, bson.M{"post_id": postId}).Decode(&post)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding post", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func UpVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var vote models.PostUpvote

		if err := c.ShouldBindJSON(&vote); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding upvote payload", "details": err.Error()})
			return
		}

		filter := bson.M{
			"user_id": vote.UserID,
			"post_id": vote.PostID,
		}

		upvote_counts, err := database.PostUpVoteCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting upvote counts", "details": err.Error()})
			return
		}
		downvote_counts, err := database.PostDownVoteCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting downvote counts", "details": err.Error()})
			return
		}
		if upvote_counts > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User has already upvote"})
			return
		}

		if downvote_counts > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User has already downvote"})
			return
		}

		result, err := database.PostUpVoteCollection.InsertOne(ctx, vote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upvoting post", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			database.PostCollection.UpdateOne(ctx, bson.M{"post_id": vote.PostID}, bson.M{"$inc": bson.M{"up_vote": 1}})
		}

		c.JSON(http.StatusCreated, result)

	}
}

func DownVotePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var vote models.PostDownvote

		if err := c.ShouldBindJSON(&vote); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding upvote payload", "details": err.Error()})
			return
		}

		filter := bson.M{
			"user_id": vote.UserID,
			"post_id": vote.PostID,
		}

		upvote_counts, err := database.PostUpVoteCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting upvote counts", "details": err.Error()})
			return
		}
		downvote_counts, err := database.PostDownVoteCollection.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting downvote counts", "details": err.Error()})
			return
		}
		if upvote_counts > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User has already upvote"})
			return
		}

		if downvote_counts > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User has already downvote"})
			return
		}

		result, err := database.PostDownVoteCollection.InsertOne(ctx, vote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error downvoting post", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			database.PostCollection.UpdateOne(ctx, bson.M{"post_id": vote.PostID}, bson.M{"$inc": bson.M{"down_vote": -1}})
		}

		c.JSON(http.StatusCreated, result)

	}
}
