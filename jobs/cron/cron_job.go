package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/EsanSamuel/Reddit_Clone/database"
	"github.com/EsanSamuel/Reddit_Clone/jobs/workers"
	"github.com/EsanSamuel/Reddit_Clone/models"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CronJob() {
	c := cron.New(cron.WithSeconds())
	defer c.Stop()

	_, err := c.AddFunc("* * * * * *", func() {
		fmt.Println("Cron job ran at ", time.Now().Format("15:04:05"))

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var posts []models.Post

		cursor, err := database.PostCollection.Find(ctx, bson.M{})

		if err != nil {
			fmt.Println(err)
			return
		}

		if err := cursor.All(ctx, &posts); err != nil {
			fmt.Println(err)
			return
		}

		for _, post := range posts {
			if time.Since(post.UpdatedAt) >= 24*time.Hour {
				fmt.Println("Post", post.PostID, "hasn't been modified in the last one day. skipping...")
				continue
			}
			workers.AISummaryQueue(post.PostID)
		}

	})
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Start()
	select {}
}
