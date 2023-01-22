package email

import (
	"encoding/json"
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"context"

	"github.com/k3a/html2text"
	"jwt-auth/initializers"
	Redis "jwt-auth/utils/redis"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
)
// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(data *EmailData,templateType string) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := data.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates/email")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	log.Println(data)
	emailTemplateFile := checkEmailType(templateType)

	template.ExecuteTemplate(&body, emailTemplateFile, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}

}

func RunSendEmailJob(data interface{},emailTemplate string) {

	D,_ := json.Marshal(data)
	Data := map[string]interface{}{ "value": string(D) , "template" : emailTemplate }

	rdq := Redis.NewClient()
	err := rdq.XAdd(context.Background(), &redis.XAddArgs{
		///this is the name we want to give to our stream
		///in our case we called it send_order_emails
		//note you can have as many stream as possible
		//such as one for email...another for notifications
		Stream:       streamName,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		//values is the data you want to send to the stream
		//in our case we send a map with email and message keys
		Values: Data,
	}).Err()
	if err != nil {
		log.Fatal( err)
		return
	}
	log.Println("email sending in background")
}

