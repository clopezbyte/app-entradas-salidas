package utils

import (
	"encoding/base64"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailWithCSV(toEmail, subject, plainTextContent string, csvData []byte) error {
	from := mail.NewEmail("App Entradas/Salidas", "noreply@yourdomain.com")
	to := mail.NewEmail("", toEmail)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, plainTextContent)

	attachment := mail.NewAttachment()
	attachment.SetContent(base64.StdEncoding.EncodeToString(csvData))
	attachment.SetType("text/csv")
	attachment.SetFilename("filtered_data.csv")
	attachment.SetDisposition("attachment")

	message.AddAttachment(attachment)

	client := sendgrid.NewSendClient("YOUR_SENDGRID_API_KEY")
	_, err := client.Send(message)
	return err
}
