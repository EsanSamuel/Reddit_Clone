package models

import "go.mongodb.org/mongo-driver/v2/bson"

type PostUpvote struct {
	ID     bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	PostID string        `json:"post_id" bson:"post_id"`
	UserID string        `json:"user_id" bson:"user_id"`
	Value  int           `json:"value" bson:"value"`
}

type PostDownvote struct {
	ID     bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	PostID string        `json:"post_id" bson:"post_id"`
	UserID string        `json:"user_id" bson:"user_id"`
	Value  int           `json:"value" bson:"value"`
}

type CommentUpvote struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	CommentID string        `json:"comment_id" bson:"comment_id"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Value     int           `json:"value" bson:"value"`
}

type CommentDownvote struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	CommentID string        `json:"comment_id" bson:"comment_id"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Value     int           `json:"value" bson:"value"`
}

/*post_votes
----------
id              UUID (PK)
user_id         UUID (FK → users.id)
post_id         UUID (FK → posts.id)
value           INT CHECK (value IN (1, -1))

UNIQUE(user_id, post_id)
*/
