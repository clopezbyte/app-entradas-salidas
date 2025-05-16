package utils

import (
	_ "image/png"
	"log"

	"fmt"

	"context"

	"github.com/clopezbyte/app-entradas-salidas/models"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func handleClientEmailNotification(ctx context.Context, entrada models.Entradas) {
	clientDoc, err := firestoreClient.Collection("customers").
		Where("name", "==", entrada.Cliente).
		Limit(1).
		Documents(ctx).
		Next()
	if err == iterator.Done {
		log.Warnf("No customer found for name: %s", entrada.Cliente)
		return
	} else if err != nil {
		log.Errorf("Failed to query Firestore: %v", err)
		return
	}

	var customer struct {
		Email string `firestore:"email"`
		Rep   string `firestore:"rep"`
	}
	if err := clientDoc.DataTo(&customer); err != nil {
		log.Errorf("Failed to map customer document: %v", err)
		return
	}

	// Construct and send email
	body, err := GenerateEmailBody(EmailData{
		CustomerName: entrada.Cliente,
		RepName:      customer.Rep,
		DeliveryType: entrada.TipoDelivery,
		DetailsURL:   fmt.Sprintf("https://your-app.com/entradas/%d", entrada.NumeroRemisionFactura),
	})
	if err != nil {
		log.Errorf("Email generation error: %v", err)
		return
	}

	if err := sendEmail(customer.Email, "Nueva recepción de mercancía", body); err != nil {
		log.Errorf("Email send error: %v", err)
	}
}
