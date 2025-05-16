package models

type Customer struct {
	Code    string `firestore:"code"`
	Email   string `firestore:"email"`
	RepName string `firestore:"rep_name"`
}
