package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateSubreddit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var subreddit models.SubReddit
		var member models.SubRedditMembers

		if err := c.ShouldBindJSON(&subreddit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error binding subreddit payload", "details": err.Error()})
			return
		}

		subreddit.CreatedAt = time.Now()
		subreddit.UpdatedAt = time.Now()
		subreddit.SubRedditId = bson.NewObjectID().Hex()

		result, err := database.SubredditCollection.InsertOne(ctx, subreddit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating subreddit", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			member.MemberId = bson.NewObjectID().Hex()
			member.JoinedAt = time.Now()
			member.UserID = subreddit.CreatorId
			member.Role = "MODERATOR"
			member.SubRedditId = subreddit.SubRedditId
			_, err := database.MemberCollection.InsertOne(ctx, member)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error adding creator to the subreddit member", "details": err.Error()})
				return
			}
			database.SubredditCollection.UpdateOne(ctx, bson.M{"subreddit_id": member.SubRedditId}, bson.M{"$inc": bson.M{"members_count": 1}})
		}

		c.JSON(http.StatusCreated, gin.H{"subreddit": subreddit, "result": result})

	}
}

func JoinSubreddit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var member models.SubRedditMembers

		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error binding member payload", "details": err.Error()})
			return
		}

		member.MemberId = bson.NewObjectID().Hex()
		member.JoinedAt = time.Now()

		member_count, err := database.MemberCollection.CountDocuments(ctx, bson.M{"user_id": member.UserID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting member", "details": err.Error()})
			return
		}

		if member_count > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "user has already joined this subreddit"})
			return
		}

		result, err := database.MemberCollection.InsertOne(ctx, member)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error joining subreddit", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			database.SubredditCollection.UpdateOne(ctx, bson.M{"subreddit_id": member.SubRedditId}, bson.M{"$inc": bson.M{"members_count": 1}})
		}

		c.JSON(http.StatusOK, member)
	}
}

func AddModerators() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var member models.SubRedditMembers

		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error binding member payload", "details": err.Error()})
			return
		}

		member.MemberId = bson.NewObjectID().Hex()
		member.JoinedAt = time.Now()
		member.Role = "MODERATOR"

		member_count, err := database.MemberCollection.CountDocuments(ctx, bson.M{"user_id": member.UserID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting member", "details": err.Error()})
			return
		}

		if member_count > 0 {
			_, err := database.MemberCollection.UpdateOne(ctx, bson.M{"user_id": member.UserID}, bson.M{"$set": bson.M{"role": "MODERATOR"}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user role to moderator", "details": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "user role has been changed to moderator"})
			return
		}

		result, err := database.MemberCollection.InsertOne(ctx, member)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error adding moderator", "details": err.Error()})
			return
		}

		if result.Acknowledged {
			database.SubredditCollection.UpdateOne(ctx, bson.M{"subreddit_id": member.SubRedditId}, bson.M{"$inc": bson.M{"members_count": 1}})
		}

		c.JSON(http.StatusOK, member)
	}
}
