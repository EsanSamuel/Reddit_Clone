package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/config"
	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
		c.JSON(http.StatusOK, summary)
	}

}
