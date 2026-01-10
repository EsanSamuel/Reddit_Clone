package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/config"
	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/gocraft/work"
	"github.com/resend/resend-go/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Context struct {
	Email    string
	UserId   string
	PostId   string
	Post     models.Post
	Comments []models.Comment
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Background job is running:", job.Name, "ID:", job.ID)
	return next()
}

func (c *Context) FindUser(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["user_id"]; ok {
		c.Email = job.ArgString("email_addr")
		c.UserId = job.ArgString("user_id")
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	return next()
}

func (c *Context) SendWelcomeEmail(job *work.Job) error {
	email := c.Email
	fmt.Println(email)
	RESEND_API_KEY := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(RESEND_API_KEY)

	params := &resend.SendEmailRequest{
		From: "Acme <noreply@mikaelsoninitiative.org>",
		To:   []string{email},
		Html: `<div style="max-width: 500px; margin: 0 auto; font-family: Arial, sans-serif; background-color: #ffffff; padding: 30px; border-radius: 8px; border: 1px solid #e5e7eb;">

  <h2 style="color: #111827; text-align: center; margin-bottom: 10px;">
    Welcome to Reddit 🎉
  </h2>

  <p style="color: #374151; font-size: 15px; text-align: center;">
    Hey there 👋
  </p>

  <p style="color: #374151; font-size: 15px; text-align: center; line-height: 1.5;">
    We’re excited to have you on <b>Reddit</b>! Your account has been successfully created, and you’re all set to start exploring communities, sharing ideas, and joining conversations that matter to you.
  </p>

  <div style="text-align: center; margin: 30px 0;">
    <a href=""
       style="
         background-color: #2563eb;
         color: #ffffff;
         padding: 14px 30px;
         text-decoration: none;
         border-radius: 6px;
         font-weight: bold;
         display: inline-block;
         font-size: 16px;
       ">
      Get Started
    </a>
  </div>

  <p style="color: #6b7280; font-size: 14px; text-align: center; line-height: 1.4;">
    If you have any questions, feel free to reply to this email — we’re happy to help.
  </p>

</div>
`,
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(sent.Id)
	return nil
}

func (c *Context) FindPost(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["post_id"]; ok {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		c.PostId = job.ArgString("post_id")

		err := database.PostCollection.FindOne(ctx, bson.M{"post_id": c.PostId}).Decode(&c.Post)

		if err != nil {
			return err
		}

		// Fetch all Comment
		cursor, err := database.CommentCollection.Find(ctx, bson.M{"post_id": c.PostId})
		if err := cursor.All(ctx, &c.Comments); err != nil {
			return err
		}
		defer cursor.Close(ctx)

		if err := job.ArgError(); err != nil {
			return err
		}
	}
	return next()
}

func (c *Context) SendAISummary(job *work.Job) error {
	commentsJSON, err := json.Marshal(c.Comments)
	post := c.Post

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
	fmt.Println(summary)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (c *Context) GeneratePostEmbeddings(job *work.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	postId := c.PostId

	var post models.Post
	var comments []models.Comment

	if err := database.PostCollection.FindOne(ctx, bson.M{"post_id": postId}).Decode(&post); err != nil {
		return err
	}

	cursor, err := database.CommentCollection.Find(ctx, bson.M{"post_id": postId})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &comments); err != nil {
		return err
	}

	var commentTexts []string
	for _, c := range comments {
		commentTexts = append(commentTexts, c.Content)
	}

	embeddingContent := post.Title + "\n" + post.Content
	if len(commentTexts) > 0 {
		embeddingContent += "\n" + strings.Join(commentTexts, "\n")
	}

	fmt.Println(embeddingContent)

	embeddings := config.AIEmbeddings(embeddingContent)
	//fmt.Println(embeddings)

	_, err = database.PostCollection.UpdateOne(
		ctx,
		bson.M{"post_id": postId},
		bson.M{"$set": bson.M{"embeddings": embeddings}},
	)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
