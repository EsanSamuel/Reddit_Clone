package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/config"
	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/helpers"
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

			var post models.Post
			var comments []models.Comment

			err := database.PostCollection.FindOne(ctx, bson.M{"post_id": postId}).Decode(&post)

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

			for i := range allChunks {
				if allChunks[i].Embedding == nil {
					allChunks[i].Embedding = config.AIEmbeddings(allChunks[i].Text)
				}
			}

			//var scores []float32
			for i, chunk := range allChunks {
				score := helpers.CosineSimilarity(queryEmbeddings, chunk.Embedding)
				allChunks[i].Score = score
			}

			sort.Slice(allChunks, func(i, j int) bool {
				return allChunks[i].Score > allChunks[j].Score
			})

			var scores []float32
			for i := range allChunks {
				scores = append(scores, allChunks[i].Score)
			}

			sort.Slice(scores, func(i, j int) bool {
				return scores[i] > scores[j]
			})

			topK := 3
			var content []string
			for i, chunk := range allChunks {
				if i < topK && chunk.Score > 0.35 {
					content = append(content, chunk.Text)
				}
			}

			contentString, err := json.Marshal(content)

			prompt := fmt.Sprintf(`You are an AI assistant. Use the following content to answer the user's query.

                    **Instructions:**
                        1. Only use the information provided in the relevant content chunks.
                        2. Provide a clear, concise, and informative answer.
                        3. Highlight disagreements, recurring ideas, or differing opinions if present.
                        4. Keep the tone neutral, factual, and professional.
                        5. Do not include information not present in the content.

                    **User Query: "%s" ** 
                    

                    **Relevant Content Chunks:"%s"** 
                        

                    **Answer:**
                     `, query, string(contentString))

			answer, err := config.Ai(prompt)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(content)

			c.JSON(http.StatusOK, gin.H{"scores": scores, "Answer": answer})
		}

	}
}
