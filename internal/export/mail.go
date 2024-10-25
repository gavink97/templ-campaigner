package export

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Request struct {
	To      []string
	From    string
	Subject string
	Body    string
	Images  []Image
}

type RequestParams struct {
	To      []string
	From    string
	Subject string
	Body    string
	Images  []Image
}

func NewRequest(params RequestParams) *Request {
	return &Request{
		To:      params.To,
		From:    params.From,
		Subject: params.Subject,
		Body:    params.Body,
		Images:  params.Images,
	}
}

// verify email addresses before sending the email
func (req *Request) SendEmail() error {
	err := godotenv.Load()
	if err != nil {
		slog.Info(fmt.Sprintf("Error loading .env file: %v", err))
		return err
	}

	privateKey := os.Getenv("PRIVATE_KEY")

	server := mail.NewSMTPClient()

	server.Host = "smtp.zeptomail.com"
	server.Port = 587
	server.Username = os.Getenv("ZEPTOMAIL_CLIENT_ID")
	server.Password = os.Getenv("ZEPTOMAIL_SECRET")
	server.Encryption = mail.EncryptionSTARTTLS
	server.Authentication = mail.AuthPlain

	// Variable to keep alive connection
	server.KeepAlive = true

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if req.To == nil {
		err = errors.New("There are no contacts in To")
		slog.Error(err.Error())
		return err
	}

	for _, to := range req.To {
		email := mail.NewMSG()

		email.SetFrom(fmt.Sprintf("Gavin Kondrath <%s>", req.From)).
			AddTo(to).
			// AddCc("otherto@example.com").
			SetSubject(req.Subject).
			SetListUnsubscribe("<mailto:unsubscribe@example.com?subject=https://example.com/unsubscribe>")

		email.SetBody(mail.TextHTML, req.Body)

		// cid
		for _, image := range req.Images {
			email.Attach(&mail.File{FilePath: image.Path, Name: image.Name, Inline: image.Inline})
		}

		// Delivery Status Notification (DSN)
		email.SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

		if privateKey != "" {
			options := dkim.NewSigOptions()
			options.PrivateKey = []byte(privateKey)
			options.Domain = "gav.ink"
			options.Selector = "default"
			options.SignatureExpireIn = 3600
			options.Headers = []string{"from", "date", "mime-version", "received", "received"}
			options.AddSignatureTimestamp = true
			options.Canonicalization = "relaxed/relaxed"

			email.SetDkim(options)
		}

		if email.Error != nil {
			slog.Error(email.Error.Error())
			return email.Error
		}

		err = email.Send(smtpClient)
		if err != nil {
			slog.Error(err.Error())
			return err
		} else {
			slog.Info("Email Sent")
		}
	}
	return nil
}
