package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/config"
	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/helpers"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var redisClient = config.Redis
var logger = config.InitLogger()

func ThreadsSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		postId := c.Param("post_id")

		// Fetch Post and Post threads
		var post models.Post
		var comments []models.Comment

		err := database.PostCollection.FindOne(ctx, bson.M{"post_id": postId}).Decode(&post)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding post", "details": err.Error()})
			return
		}

		// Fetch all Comment
		cursor, err := database.CommentCollection.Find(ctx, bson.M{"post_id": postId})
		if err := cursor.All(ctx, &comments); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comments", "details": err.Error()})
			return
		}

		commentsJSON, err := json.Marshal(comments)

		prompt := fmt.Sprintf(`You are an AI assistant. I will provide you with a post and its associated comments. Summarize the content for a user in a concise and informative way. Include the following:

                 1. **Post Summary**: What is the post about? Key points only.
                 2. **Comment Thread Summary**: Main opinions, arguments, or insights from the comments.
                 3. **Tone**: Neutral and clear.
                 4. **Optional**: Highlight if there are disagreements or recurring ideas.

                    Here is the data:

                    Post:
                    Title: "%s"
                    Content: "%s"
                    Type: "%s"  // e.g., text, image, link
                    Tags: %v
                    Comments: %s
                    `, post.Title, post.Content, post.Type, post.Tags, string(commentsJSON))

		summary, err := config.Ai(prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error summarizing posts", "details": err.Error()})
		}

		logger.INFO(summary)
		logger.DEBUG("Debug detail")
		c.JSON(http.StatusOK, summary)
	}

}

func SeachPostDetailsWithAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		postId := c.Param("postId")
		query := c.Query("query")

		embedding := config.AIEmbeddings(query)
		if embedding != nil {
			queryEmbeddings := embedding
			fmt.Println(queryEmbeddings)

			cmd := redisClient.Get(ctx, postId)

			data, err := cmd.Bytes()
			if err != nil {
				if err == redis.Nil {
					log.Println("No chunks in Redis, continuing normal logic...")
				} else {
					log.Println("Redis error:", err)
					return
				}
			}

			if len(data) > 0 {
				fmt.Println("Chunks found in Redis")

				var allChunks []helpers.Chunk
				if err := json.Unmarshal(data, &allChunks); err != nil {
					log.Println("Failed to unmarshal redis chunks:", err)
					return

				}
				scores, answer := helpers.ProcessChunks(allChunks, queryEmbeddings, query)
				logger.INFO(answer)
				c.JSON(http.StatusOK, gin.H{"scores": scores, "Answer": answer})
			} else {

				var post models.Post
				var comments []models.Comment

				err = database.PostCollection.FindOne(ctx, bson.M{"post_id": postId}).Decode(&post)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding post", "details": err.Error()})
					return
				}

				cursor, err := database.CommentCollection.Find(ctx, bson.M{"post_id": postId})
				if err != nil {
					return
				}
				defer cursor.Close(ctx)

				if err := cursor.All(ctx, &comments); err != nil {
					fmt.Println(err)
					return
				}

				titleChunk := helpers.ChunkText(postId, "title", post.Title)
				//fmt.Println(titleChunk)

				var commentsChunks []helpers.Chunk
				for _, c := range comments {
					chunks := helpers.ChunkText(postId, "comments", c.Content)
					commentsChunks = append(commentsChunks, chunks...)
				}
				allChunks := append(titleChunk, commentsChunks...)
				//fmt.Println(commentsChunks)

				//result := helpers.CosineSimilarity(post.Embeddings, queryEmbeddings)

				scores, answer := helpers.ProcessChunks(allChunks, queryEmbeddings, query)

				data, err = json.Marshal(allChunks)
				if err != nil {
					fmt.Println("JSON marshal error:", err.Error())
					return
				}

				if err := redisClient.Set(ctx, post.PostID, data, 86400*time.Second); err != nil {
					fmt.Println("Error storing chunks in redis", err)
				}

				logger.INFO(answer)
				c.JSON(http.StatusOK, gin.H{"scores": scores, "Answer": answer})
			}

		}
	}
}
