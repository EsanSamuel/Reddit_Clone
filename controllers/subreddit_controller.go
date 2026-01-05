package controllers

import (
	"context"
	"fmt"
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

func GetSubReddit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var subreddits []models.SubReddit

		filter := bson.M{}
		findOptions := options.Find()

		// Search Subreddit
		if s := strings.TrimSpace(c.Query("search")); s != "" {
			safe := regexp.QuoteMeta(s)
			filter = bson.M{
				"$or": []bson.M{
					{
						"name": bson.M{
							"$regex":   safe,
							"$options": "i",
						},
					},
					{
						"description": bson.M{
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

		cursor, err := database.SubredditCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding subreddits", "details": err.Error()})
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &subreddits); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding subreddits", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, subreddits)

	}
}

func GetSubRedditUserJoined() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var memberships []models.SubRedditMembers
		var subreddit []models.SubReddit

		user_id := c.Param("user_id")

		findOptions := options.Find()

		// Sort Subreddit
		if sort := strings.TrimSpace(c.Query("sort")); sort != "" {
			switch sort {
			case "asc":
				findOptions.SetSort(bson.D{{"joined_at", 1}})
			case "desc":
				findOptions.SetSort(bson.D{{"joined_at", -1}})
			}

		}

		cursor, err := database.MemberCollection.Find(ctx, bson.M{"user_id": user_id}, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding member", "details": err.Error()})
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &memberships); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding subreddits", "details": err.Error()})
			return
		}
		fmt.Println("No of subreddit joined:", len(memberships))

		var subredditUserJoined []string
		for _, membership := range memberships {
			fmt.Println("Subreddit ID: ", membership.SubRedditId)
			subredditUserJoined = append(subredditUserJoined, membership.SubRedditId)
		}
		fmt.Println("No of subreddit joined:", len(subredditUserJoined))

		for _, subredditId := range subredditUserJoined {
			filter := bson.M{"subreddit_id": subredditId}
			if s := strings.TrimSpace(c.Query("search")); s != "" {
				safe := regexp.QuoteMeta(s)

				filter["$and"] = []bson.M{
					{
						"user_id": user_id,
					},
					{
						"$or": []bson.M{
							{
								"name": bson.M{
									"$regex":   safe,
									"$options": "i",
								},
							},
							{
								"description": bson.M{
									"$regex":   safe,
									"$options": "i",
								},
							},
						},
					},
				}
			}

			result, err := database.SubredditCollection.Find(ctx, filter)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding subreddit", "details": err.Error()})
				return
			}

			defer result.Close(ctx)

			if err := result.All(ctx, &subreddit); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding subreddits", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, subreddit)
		}
	}
}

func GetSubRedditById() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		subredditId := c.Param("id")

		var subreddit models.SubReddit

		err := database.SubredditCollection.FindOne(ctx, bson.M{"subreddit_id": subredditId}).Decode(&subreddit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding subreddit", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, subreddit)
	}
}
