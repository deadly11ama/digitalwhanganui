package main

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/gomail.v2"
)

var emailTemplates *template.Template

func administratorEmail() string {
	return "DigitalWhanganui <" + Config.AdminEmailAddress + ">"
}

func fromEmail() string {
	return "DigitalWhanganui <" + Config.FromEmailAddress + ">"
}

func sendMail(to, subject, template string, data map[string]string) {
	checkHeaderText(to)
	checkHeaderText(subject)

	var b bytes.Buffer

	err := emailTemplates.ExecuteTemplate(&b, template, data)
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", b.String())

	d := gomail.NewDialer(Config.SMTPServer, Config.SMTPPort, Config.SMTPUser, Config.SMTPPassword)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

func (l *Listing) FullAdminEmail() string {
	return l.AdminFirstName + " " + l.AdminLastName + " <" + l.AdminEmail + ">"
}

// Ensure headerText is not being used to inject additional email headers
func checkHeaderText(headerText string) {
	if strings.ContainsAny(headerText, "\n\r") {
		panic("Possible Email header injection: " + headerText)
	}
}

func emailErrorMsg(msg string, logger *log.Logger) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("Panic when emailing error message:", err)
		}
	}()

	sendMail(Config.ErrorEmailAddress, "Digital Whanganui Error", "error.tmpl", map[string]string{"error": msg})
}

func initEmail() {
	pattern := filepath.Join(Config.EmailTemplateDir, "*.tmpl")
	emailTemplates = template.Must(template.ParseGlob(pattern))

}
