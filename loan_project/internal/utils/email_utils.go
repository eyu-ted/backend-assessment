package utils

import (
	"log"

	"gopkg.in/mail.v2"
)

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	SenderEmail string
	SenderName  string
	SenderPass  string
}

func SendVerificationEmail(emailConfig EmailConfig, recipientEmail, token string) error {
	mailer := mail.NewMessage()

	mailer.SetHeader("From", emailConfig.SenderName+" <"+emailConfig.SenderEmail+">")

	mailer.SetHeader("To", recipientEmail)

	// Set the subject
	mailer.SetHeader("Subject", "Email Verification")

	// Set the body
	body := "Please verify your email by clicking the link below:\n\n"
	body += "http://localhost:8080/users/verify-email?token=" + token + "&email=" + recipientEmail
	mailer.SetBody("text/plain", body)

	// Configure the SMTP server
	dialer := mail.NewDialer(emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.SenderEmail, emailConfig.SenderPass)

	// Send the email
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Verification email sent to %s", recipientEmail)
	return nil
}

func SendPasswordResetEmail(emailConfig EmailConfig, recipientEmail, token string) error {
	mailer := mail.NewMessage()

	
	mailer.SetHeader("From", emailConfig.SenderName+" <"+emailConfig.SenderEmail+">")

	mailer.SetHeader("To", recipientEmail)

	mailer.SetHeader("Subject", "Password Reset Request")


	body := "Click the link below to reset your password:\n\n"
	body += "http://localhost:8080/users/password-update?token=" + token + "&email=" + recipientEmail
	mailer.SetBody("text/plain", body)

	dialer := mail.NewDialer(emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.SenderEmail, emailConfig.SenderPass)


	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Password reset email sent to %s", recipientEmail)
	return nil
}
