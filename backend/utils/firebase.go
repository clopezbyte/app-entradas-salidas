package utils

import (
	"context"
	"errors"

	"bytes"
	"log"
	"text/template"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/clopezbyte/app-entradas-salidas/models"
	"google.golang.org/api/iterator"
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

//GenerateEmailBody generates the email body for the notification

func GenerateEmailBody(data EmailData) (string, error) {
	const tpl = `
		Hola {{.RepName}},

		Se ha registrado una nueva {{.TipoDelivery}} para el cliente "{{.Cliente}}".

		Fecha de entrada: {{.FechaRecepcion}}
		Bodega: {{.BodegaRecepcion}}
		Cantidad: {{.Cantidad}}
		Numero de remisión: {{.NumeroRemision}}
		Con proveedor: {{.ProveedorRecepcion}}
		Evidencia de entrada: {{.EvidenciaRecepcion}}

		Saludos,
		Buho Logistics
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

// ClientNotification
func HandleClientEmailNotification(ctx context.Context, firestoreClient *firestore.Client, entrada models.EntradasData) {
	clientDoc, err := firestoreClient.Collection("customers").
		Where("name", "==", entrada.Cliente).
		Limit(1).
		Documents(ctx).
		Next()
	if err == iterator.Done {
		log.Printf("No customer found for name: %s", entrada.Cliente)
		return
	} else if err != nil {
		log.Printf("Failed to query Firestore: %v", err)
		return
	}

	var customer struct {
		Email   string `firestore:"email"`
		RepName string `firestore:"rep_name"`
	}
	if err := clientDoc.DataTo(&customer); err != nil {
		log.Printf("Failed to map customer document: %v", err)
		return
	}

	// Construct and send email
	body, err := GenerateEmailBody(EmailData{
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
		log.Printf("Email generation error: %v", err)
		return
	}

	if err := sendEmail(customer.Email, "Nueva recepción de mercancía", body); err != nil {
		log.Printf("Email send error: %v", err)
	}
}
