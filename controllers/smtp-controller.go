package controllers

import (
	"fmt"
	"net/smtp"
)

type SmtpServer struct {
	Host string
	Port string
}


const Email = "travelohiVR@gmail.com"
const EmailPassword = "mzfd xssk evpa oozk"
var SMTP_SERVER = SmtpServer{Host: "smtp.gmail.com", Port: "587"}

func SendEmail(server *SmtpServer, from string, to []string, password string, subject string, content string) error {
	formattedSubject := "Subject: [TraveloHI] " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><h2>%s</h2>%s</body></html>", subject, content)
	msg := []byte(formattedSubject + mime + body)
	auth := smtp.PlainAuth("", from, password, server.Host)

	err := smtp.SendMail(server.Host + ":" + server.Port, auth, from, to, []byte(msg))

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil

}