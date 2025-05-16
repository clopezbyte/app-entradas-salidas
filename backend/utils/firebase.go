package utils

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/clopezbyte/app-entradas-salidas/models"
	"github.com/mailersend/mailersend-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Initializes the Firebase Admin SDK and returns the Auth client
func InitializeFirebase() (*auth.Client, error) {
	// Initialize the Firebase app with default credentials
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Get an Auth client from the Firebase app
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Verifies the Firebase ID token from the request header
func VerifyIDToken(idToken string) (*auth.Token, error) {
	client, err := InitializeFirebase()
	if err != nil {
		return nil, err
	}

	// Verify the ID token
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return token, nil
}

// Extract token from header
func GetTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("missing Authorization token")
	}

	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid token format")
	}

	return authHeader[7:], nil
}

// generateEmailBody generates the email body using Go templates.
func generateEmailBody(data models.EmailData) (string, error) {
	const tpl = `
Hola,

Se ha registrado una nueva devolución para el cliente "{{.Cliente}}".

Fecha de entrada: {{.FechaRecepcion}}
Bodega: {{.BodegaRecepcion}}
Cantidad: {{.Cantidad}}
Numero de remisión: {{.NumeroRemision}}
Con proveedor: {{.ProveedorRecepcion}}
Link a evidencia de entrada: {{.EvidenciaRecepcion}}

Saludos,  
Buho Logistics  

(Correo automático, favor de no responder.)
`
	t, err := template.New("email").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// sendEmail sends an email using MailerSend API.
func sendEmail(to string, repName string, subject, body string) error {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	if apiKey == "" {
		return errors.New("MAILERSEND_API_KEY not set in environment")
	}

	ms := mailersend.NewMailersend(apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Initialize the message with the builder
	message := ms.Email.NewMessage()

	message.SetFrom(mailersend.From{
		Email: "buhologistics@test-68zxl27957m4j905.mlsender.net", // must match your MailerSend domain
		Name:  "Buho Logistics",
	})

	message.SetRecipients([]mailersend.Recipient{
		{
			Email: to,
			Name:  repName,
		},
	})

	message.SetSubject(subject)
	message.SetText(body)
	message.SetHTML(body)

	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		log.Printf("Mailersend send failed: %v", err)
		return err
	}

	log.Printf("X-Message-Id: %s", res.Header.Get("X-Message-Id")) // helpful for debugging	return nil
	log.Printf("Email sent successfully to %s (%s)", repName, to)
	return nil
}

// HandleClientEmailNotification looks up the customer and sends an email notification.
func HandleClientEmailNotification(ctx context.Context, firestoreClient *firestore.Client, entrada models.Entradas) {
	log.Printf("Looking up customer with ID: %s", entrada.Cliente)

	// Fetch customer document by ID
	docSnap, err := firestoreClient.Collection("customers").Doc(entrada.Cliente).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Document not found, handle accordingly
			log.Printf("No customer found with ID: %s", entrada.Cliente)
			return
		}
		// Other errors
		log.Printf("Firestore error: %v", err)
		return
	}

	var customer models.Customer
	if err := docSnap.DataTo(&customer); err != nil {
		log.Printf("Failed to parse customer document: %v", err)
		return
	}

	// Build the email content
	body, err := generateEmailBody(models.EmailData{
		Email:              customer.Email,
		RepName:            customer.RepName,
		BodegaRecepcion:    entrada.BodegaRecepcion,
		Cantidad:           int(entrada.Cantidad),
		Comentarios:        entrada.Comentarios,
		EvidenciaRecepcion: entrada.EvidenciaRecepcion,
		FechaRecepcion:     entrada.FechaRecepcion,
		NumeroRemision:     int(entrada.NumeroRemisionFactura),
		PersonaRecepcion:   entrada.PersonaRecepcion,
		ProveedorRecepcion: entrada.ProveedorRecepcion,
		Cliente:            entrada.Cliente,
		TipoDelivery:       entrada.TipoDelivery,
	})
	if err != nil {
		log.Printf("Email body generation failed: %v", err)
		return
	}

	// Send the email
	if err := sendEmail(customer.Email, customer.RepName, "Nueva devolución de mercancía", body); err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
