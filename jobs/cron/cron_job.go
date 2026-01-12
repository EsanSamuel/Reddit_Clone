package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func CronJob() {
	c := cron.New(cron.WithSeconds())
	defer c.Stop()

	_, err := c.AddFunc("@daily", func() {
		fmt.Println("Cron job ran at ", time.Now().Format("15:04:05"))

	})
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Start()
	select {}
}
