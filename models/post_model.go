package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Post struct {
	ID           bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	PostID       string        `json:"post_id" bson:"post_id"`
	Title        string        `json:"title" form:"title" bson:"title" validate:"required"`
	Content      string        `json:"content" form:"content" bson:"content" validate:"required"`
	Type         string        `json:"type" form:"type" bson:"type"`
	FileUrls     []string      `json:"file_urls" bson:"file_urls"`
	AuthorID     string        `json:"author_url" form:"author_url" bson:"author_url"`
	SubredditID  string        `json:"subreddit_id" form:"subreddit_id" bson:"subreddit_id"`
	Score        int           `json:"score" bson:"score"`
	Tags         []string      `json:"tags" form:"tags" bson:"tags"`
	CommentCount int           `json:"comment_count" bson:"comment_count"`
	UpVote       int           `json:"up_vote" bson:"up_vote"`
	DownVote     int           `json:"down_vote" bson:"down_vote"`
	Embeddings   []float32     `json:"embeddings" bson:"embeddings"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type TagsPosts struct {
	Posts     []Post `json:"posts" bson:"posts"`
	PostCount int    `json:"post_count" bson:"post_count"`
}

/*id              UUID (PK)
title           VARCHAR
content         TEXT
type            ENUM('text', 'link', 'image')
image_url       TEXT
author_id       UUID (FK → users.id)
subreddit_id    UUID (FK → subreddits.id)
score           INT DEFAULT 0
created_at      TIMESTAMP
updated_at      TIMESTAMP*/
