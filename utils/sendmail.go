package utils

import (
	"net/smtp"
)

const (
	smtpServer     = "smtp.gmail.com"
	smtpPort       = "587"
	senderEmail    = "libsongo@gmail.com"
	senderPassword = "adihcoyyaiofjysw" // App password (no spaces)
)

func SendEmail(to, subject, body string) error {
	msg := "From: " + senderEmail + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	return smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{to}, []byte(msg))
}
