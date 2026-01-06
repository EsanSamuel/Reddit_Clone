package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Cannot find .env file")
	}

	mongoDB_URI := os.Getenv("DATEBASE_URI")

	if mongoDB_URI == "" {
		fmt.Println("MongoDB URI is empty")
	}

	clientOptions := options.Client().ApplyURI(mongoDB_URI)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		fmt.Println("Error connecting to mongodb database")
	}

	return client
}

var Client *mongo.Client = Connect()

func Collection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Cannot find .env file")
	}

	DatabaseName := os.Getenv("DATABASE_NAME")

	if DatabaseName == "" {
		fmt.Println("Database name is empty")
	}

	collection := Client.Database(DatabaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}

var UserCollection *mongo.Collection = Collection("users")
var SubredditCollection *mongo.Collection = Collection("subreddit")
var MemberCollection *mongo.Collection = Collection("subreddit_members")
var PostCollection *mongo.Collection = Collection("posts")
var CommentCollection *mongo.Collection = Collection("comments")
var PostUpVoteCollection *mongo.Collection = Collection("post_upvote")
var PostDownVoteCollection *mongo.Collection = Collection("post_downvote")
