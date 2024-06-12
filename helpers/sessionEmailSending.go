package helpers

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
)

func sessionEmailSending(emailId string, mobile string, name string, city string, service string) error {
	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Email details
	from := "leadataprlf@gmail.com"
	password := "csnm smmk umvn mdow"
	to := emailId

	// Create the HTML template
	htmlTemplate := `
	<html>
	<body>
	<div class="email-header">
		<h1>Welcome to PhysioRelief</h1>
	</div>
	<h2>Booking Session Details</h2>
	<p>Email ID : %s</p>
	<p>Name : %s</p>
	<p>Mobile No : %s</p>
	<p>City : %s</p>
	<p>Service : %s</p><br/>
	<p>Best regards,
	<h5>PhysioRelief</h5>
	</p>
	</body>
	</html>
	`

	// // Format the data into the HTML template
	htmlContent := fmt.Sprintf(htmlTemplate, emailId, name, mobile, city, service)

	// Create a new multipart message
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// Set the boundary
	boundary := mw.Boundary()

	// Create the email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = "Booking Session Details"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/alternative; boundary="` + boundary + `"`

	// Write the headers to the buffer
	for k, v := range headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}
	buf.WriteString("\r\n")

	// Create the HTML part
	part, err := mw.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=utf-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		log.Fatal(err)
	}
	part.Write([]byte(htmlContent))

	// Close the multipart writer
	if err := mw.Close(); err != nil {
		log.Fatal(err)
	}

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, buf.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Email sent successfully to User!")

	to1 := "enquiryphysiorelief@gmail.com"

	// Create the HTML template
	htmlTemplate1 := `
	<html>
	<body>
	<div class="email-header">
		<h1>%s Session Booking Details</h1>
	</div>
	<p>Email ID : %s</p>
	<p>Name : %s</p>
	<p>Mobile No : %s</p>
	<p>City : %s</p>
	<p>Service : %s</p><br/>
	<p>Best regards,
	<h5>PhysioRelief</h5>
	</p>
	</body>
	</html>
	`

	// // Format the data into the HTML template
	htmlContent1 := fmt.Sprintf(htmlTemplate1, name, emailId, name, mobile, city, service)

	// Create a new multipart message
	var buf1 bytes.Buffer
	mw1 := multipart.NewWriter(&buf1)

	// Set the boundary
	boundary1 := mw1.Boundary()

	// Create the email headers
	headers1 := make(map[string]string)
	headers1["From"] = from
	headers1["To"] = to1
	headers1["Subject"] = "Booking Session Details"
	headers1["MIME-Version"] = "1.0"
	headers1["Content-Type"] = `multipart/alternative; boundary="` + boundary1 + `"`

	// Write the headers to the buffer
	for k, v := range headers1 {
		buf1.WriteString(k + ": " + v + "\r\n")
	}
	buf1.WriteString("\r\n")

	// Create the HTML part
	part1, err := mw1.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=utf-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		log.Fatal(err)
	}
	part1.Write([]byte(htmlContent1))

	// Close the multipart writer
	if err := mw1.Close(); err != nil {
		log.Fatal(err)
	}

	// Authentication
	auth1 := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth1, from, []string{to1}, buf1.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Email sent successfully to Our User!")
	return nil
}
