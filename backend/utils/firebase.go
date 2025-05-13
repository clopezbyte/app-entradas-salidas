package utils

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
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
