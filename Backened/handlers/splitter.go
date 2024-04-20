package handlers

import (
	"encoding/json"
	"net/http"
	"net/smtp"
)

// EmailRequest represents the structure of the incoming JSON request
type EmailRequest struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ContentType string   `json:"contentType"`
}

// SendEmails handles the API request for sending emails
func SendEmails(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into EmailRequest struct
	var req EmailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Compose the email message
	msg := []byte("To: " + req.To[0] + "\r\n" +
		"Subject: " + req.Subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: " + req.ContentType + "\r\n" +
		"\r\n" +
		req.Body)

	// Set up authentication credentials
	auth := smtp.PlainAuth("", "bookmycinemaapp@gmail.com", "edrb zutv apyo ewjw", "smtp.gmail.com")

	// Send the email with authentication
	err = smtp.SendMail("smtp.gmail.com:587", auth, "bookmycinemaapp@gmail.com", req.To, msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
