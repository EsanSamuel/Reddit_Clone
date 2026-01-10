package workers

import (
	"fmt"
	"log"

	"github.com/EsanSamuel/Reddit_Clone/jobs"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Redis connection
func NewRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

var redisPool *redis.Pool = NewRedisPool(":6379")

func AISummaryQueue(postId string) {
	var enqueuer = work.NewEnqueuer("ai_summaryQueue", redisPool)

	_, err := enqueuer.Enqueue("send_ai_summary", work.Q{"post_id": postId})

	if err != nil {
		fmt.Println("Error queuing ai summary", err.Error())
		log.Fatal(err)
	}
}

func AISummaryWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "ai_summaryQueue", redisPool)

	worker.Middleware((*jobs.Context).Log)
	worker.Middleware((*jobs.Context).FindPost)

	worker.Job("send_ai_summary", (*jobs.Context).SendAISummary)

	worker.Start()
}

func StopAISummaryWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "ai_summaryQueue", redisPool)
	worker.Stop()
}

func SendEmailQueue(email string, userId string) {
	var enqueuer = work.NewEnqueuer("emailQueue", redisPool)

	_, err := enqueuer.Enqueue("send_welcome_email", work.Q{"email_addr": email, "user_id": userId})

	if err != nil {
		fmt.Println("Error queuing email", err.Error())
		log.Fatal(err)
	}
}

func EmailWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "emailQueue", redisPool)

	worker.Middleware((*jobs.Context).Log)
	worker.Middleware((*jobs.Context).FindUser)

	worker.Job("send_welcome_email", (*jobs.Context).SendWelcomeEmail)

	worker.Start()
}

func StopEmailWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "emailQueue", redisPool)
	worker.Stop()
}

func AIEmbeddingQueue(postId string) {
	var enqueuer = work.NewEnqueuer("ai_embeddings_queue", redisPool)

	_, err := enqueuer.Enqueue("generate_ai_embeddings", work.Q{"post_id": postId})

	if err != nil {
		fmt.Println("Error queuing ai embedding", err.Error())
		log.Fatal(err)
	}
}

func AIEmbeddingWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "ai_embeddings_queue", redisPool)

	worker.Middleware((*jobs.Context).Log)
	worker.Middleware((*jobs.Context).FindPost)

	worker.Job("generate_ai_embeddings", (*jobs.Context).GeneratePostEmbeddings)

	worker.Start()
}

func StopAIEmbeddingWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "ai_embeddings_queue", redisPool)
	worker.Stop()
}
