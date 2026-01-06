package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	ID           bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	CommentID    string        `json:"comment_id" bson:"comment_id"`
	PostID       string        `json:"post_id" bson:"post_id"`
	ParentID     string        `json:"parent_id" bson:"parent_id"`
	Content      string        `json:"content" bson:"content" validate:"required"`
	Type         string        `json:"type" bson:"type"`
	AuthorID     string        `json:"author_url" bson:"author_url"`
	Score        int           `json:"score" bson:"score"`
	CommentCount int           `json:"comment_count" bson:"comment_count"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

/*comments
--------
id              UUID (PK)
content         TEXT
author_id       UUID (FK → users.id)
post_id         UUID (FK → posts.id)
parent_id       UUID (FK → comments.id, NULLABLE)
score           INT DEFAULT 0
created_at      TIMESTAMP
updated_at      TIMESTAMP
*/
