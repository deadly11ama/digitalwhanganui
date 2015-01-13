package main

import (
	"bytes"
	"github.com/richardcrichardc/martini"
	"gopkg.in/alexcesaro/quotedprintable.v1"
	"log"
	"net/smtp"
	"path/filepath"
	"text/template"
)

var shortFromEmail = "digitalwhanganui@digitalwhanganui.org.nz"
var fromEmail = "DigitalWhanganui <" + shortFromEmail + ">"

var server = "bakerloo.richardc.net:587"
var auth = smtp.CRAMMD5Auth("digitalwhanganui@digitalwhanganui.org.nz", "apodacaGritz")

var emailTemplates *template.Template

func shortAdministratorEmail() string {
	if martini.Env == martini.Dev {
		return "richardc+digitalwhanganui@richardc.net"
	} else {
		return "digitalwhanganui@digitalwhanganui.org.nz"
	}
}

func administratorEmail() string {
	return "DigitalWhanganui <" + shortAdministratorEmail() + ">"
}

func sendMail(to, subject, template string, data map[string]string) {
	checkHeader(to)
	checkHeader(subject)

	var buf bytes.Buffer
	buf.WriteString("From: " + fromEmail + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	buf.WriteString("\r\n")

	QPWriter := quotedprintable.NewEncoder(&buf)

	err := emailTemplates.ExecuteTemplate(QPWriter, template, data)
	if err != nil {
		panic(err)
	}
	msg := buf.Bytes()

	// TODO Fix email args - almost certainly wrong but works
	err = smtp.SendMail(server, auth, fromEmail, []string{to}, msg)
	if err != nil {
		panic(err)
	}

}

func (l *Listing) FullAdminEmail() string {
	return l.AdminFirstName + " " + l.AdminLastName + " <" + l.AdminEmail + ">"
}

func checkHeader(header string) {
	// TODO
}

func emailErrorMsg(msg string, logger *log.Logger) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("Panic when emailing error message:", err)
		}
	}()

	sendMail("richardc+digitalwhanganui@richardc.net", "Digital Whanganui Error", "error.tmpl", map[string]string{"error": msg})
}

func init() {
	var dir string
	if fileExists("email-templates") {
		dir = "email-templates"
	} else {
		dir = "/usr/local/share/digitalwhanganui/email-templates"
	}

	pattern := filepath.Join(dir, "*.tmpl")
	emailTemplates = template.Must(template.ParseGlob(pattern))

}
