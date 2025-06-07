package mailer

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"os"
	"time"

	"github.com/RickinShah/BuzzChat/internal/jsonlog"
	"github.com/go-mail/mail/v2"
	"github.com/redis/go-redis/v9"
)

type EmailJob struct {
	Recipient string
	Template  string
	Data      any
}

//go:embed templates/*
var templateFS embed.FS

var logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(plainBody, "plainBody", data); err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data); err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for range 3 {
		if err = m.dialer.DialAndSend(msg); err == nil {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
	return err
}

func (m *Mailer) SendUsingConn(conn mail.SendCloser, job EmailJob) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+job.Template)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", job.Data); err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(plainBody, "plainBody", job.Data); err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(htmlBody, "htmlBody", job.Data); err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", job.Recipient)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	return mail.Send(conn, msg)
}

func EnqueueEmail(rdb *redis.Client, job EmailJob) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return rdb.RPush(context.Background(), "email_queue", payload).Err()
}
