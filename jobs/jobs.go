package jobs

import (
	"fmt"
	"os"

	"github.com/gocraft/work"
	"github.com/resend/resend-go/v3"
)

type Context struct {
	Email  string
	UserId string
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
