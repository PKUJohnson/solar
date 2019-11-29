package toolkit

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

type SendEmail struct {
	ToEmail []string
	Subject string
	Body    string
}

func (self *SendEmail) Start() {
	user := "no-reply@github.com"
	password := "xxx"
	host := "smtp.exmail.qq.com:25"

	to_stirng := strings.Join(self.ToEmail, ";")
	err := sendToMail(user, password, host, to_stirng, self.Subject, self.Body, "text")
	if err != nil {
		fmt.Println(self.ToEmail, "Send mail error!", err)
	} else {
		fmt.Println("Send mail success!", self.ToEmail)
	}
}
