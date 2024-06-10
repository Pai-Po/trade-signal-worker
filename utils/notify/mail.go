package notify

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/matcornic/hermes/v2"
	"github.com/wneessen/go-mail"
)

func sendMail(to string, subject string, body string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	portStr := os.Getenv("PORT")

	if smtpServer == "" || portStr == "" || username == "" || password == "" {
		return
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatalf("failed to parse email port: %s", err)
	}

	m := mail.NewMsg()
	if err := m.From(username); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To(to); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)
	c, err := mail.NewClient(smtpServer, mail.WithPort(int(port)), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}

func SendTestEMail(to string) {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "TradeSignal",
			Link: "https://trade-signal.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	//utils.InitTask()
	email := hermes.Email{
		Body: hermes.Body{
			Name: to,
			Intros: []string{
				"Welcome to TradeSignal! We're very excited to have you on board.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Firstname", Value: "Jon"},
				{Key: "Lastname", Value: "Snow"},
				{Key: "Birthday", Value: "01/01/283"},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with TradeSignal, please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	sendMail(to, "Welcome to TradeSignal", emailBody)
}
