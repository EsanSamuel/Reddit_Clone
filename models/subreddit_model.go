package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SubReddit struct {
	ID           bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SubRedditId  string        `json:"subreddit_id" bson:"subreddit_id"`
	Name         string        `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Description  string        `json:"description" bson:"description" validate:"required,min=2,max=200"`
	CreatorId    string        `json:"creator_id" bson:"creator_id" validate:"required"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
	MembersCount int           `json:"members_count" bson:"members_count"`
	PostsCount   int           `json:"posts_count" bson:"posts_count"`
}

type SubRedditMembers struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MemberId    string        `json:"member_id" bson:"member_id"`
	UserID      string        `json:"user_id" bson:"user_id" validate:"required"`
	SubRedditId string        `json:"subreddit_id" bson:"subreddit_id" validate:"required"`
	Role        string        `json:"role" bson:"role" validate:"oneof MEMBER MODERATOR"`
	JoinedAt    time.Time     `json:"joined_at" bson:"joined_at"`
}
