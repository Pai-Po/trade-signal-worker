package notify

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailersend/mailersend-go"
)

func sendMail(from string, fromEmail string, to string, toEmail string, subject string, bodyTemp string, bodyTempPara map[string]interface{}) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("MAILSERVER_APIKEY")
	ms := mailersend.NewMailersend(apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mailFrom := mailersend.From{
		Name:  from,
		Email: fromEmail,
	}

	mailTo := []mailersend.Recipient{
		{
			Name:  to,
			Email: toEmail,
		},
	}

	personalization := []mailersend.Personalization{
		{
			Email: toEmail,
			Data:  bodyTempPara,
		},
	}

	// tags := []string{"signal", "notify"}

	message := ms.Email.NewMessage()
	message.SetFrom(mailFrom)
	message.SetRecipients(mailTo)
	message.SetSubject(subject)
	message.SetTemplateID(bodyTemp)
	// message.SetInReplyTo("client-id")
	message.SetPersonalization(personalization)

	ms.Email.Send(ctx, message)
}

func SendWelcomeEMail(to string, userName string, confirmURL string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tempWelcomeID := os.Getenv("TEMP_WELCOME_ID")
	sendName := os.Getenv("SENDER_NAME")
	sendEmail := os.Getenv("SENDER_EMAIL")
	welcomPara := map[string]interface{}{
		"name": userName,
	}
	sendMail(sendName, sendEmail, userName, to, "Welcome to TradeSignal", tempWelcomeID, welcomPara)

}

func SendNewSignalEMail(to string, userName string, symbol string, time string, signal string, stragety string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TEMP_SIGNAL_ID := os.Getenv("TEMP_WELCOME_ID")
	sendName := os.Getenv("SENDER_NAME")
	sendEmail := os.Getenv("SENDER_EMAIL")
	sigPara := map[string]interface{}{
		"symbol":   symbol,
		"time":     time,
		"signal":   signal,
		"strategy": stragety,
	}
	sendMail(sendName, sendEmail, userName, to, "New Signal from TradeSignal", TEMP_SIGNAL_ID, sigPara)
}
