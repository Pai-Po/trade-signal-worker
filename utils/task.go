package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"trade-signal-worker/utils/notify"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

// A list of task types.
const (
	TypeWelcomeEmail = "email:welcome"
	TypeSignalEmail  = "email:signal"
)

// Task payload for any email related tasks.
type emailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
	return nil
}

func HandleSignalEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Send Signal Email to User %d", p.UserID)
	notify.SendTestEMail("pai.po.sec@gmail.com")
	return nil
}

func InitTask() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")

	connOpt := asynq.RedisClientOpt{
		Addr:     addr,
		Password: password,
	}
	srv := asynq.NewServer(
		connOpt,
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeWelcomeEmail, HandleWelcomeEmailTask)
	mux.HandleFunc(TypeSignalEmail, HandleSignalEmailTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
