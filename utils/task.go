package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
	"trade-signal-worker/utils/db"
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
	TaskID    int
	EventTime int64
	Signal    string
	Strategy  string
}

type welcomeEmailPayload struct {
	UserName   string
	UserEmail  string
	ConfirmURL string
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p welcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	// verify param
	if p.UserName == "" || p.UserEmail == "" || p.ConfirmURL == "" {
		return nil
	}

	log.Printf(" [*] Send Welcome Email to User %s %s ", p.UserName, p.UserEmail)
	notify.SendWelcomeEMail(p.UserEmail, p.UserName, p.ConfirmURL)
	return nil
}

func HandleSignalEmailTask(ctx context.Context, t *asynq.Task) error {
	var p emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	// verify param
	if p.TaskID == 0 || p.EventTime == 0 || p.Signal == "" || p.Strategy == "" {
		return nil
	}

	log.Printf(" [*] Send Signal Email to Task %d", p.TaskID)
	// get task info
	tasDB := db.NewPostgresTaskDB("")
	defer tasDB.Close()
	task, err := tasDB.GetTaskByTaskID(p.TaskID)
	if err != nil {
		log.Fatalf("Error getting task by taskID: %v", err)
		return err
	}
	userDB, err := db.NewPostgresUserDB("")
	defer userDB.Close()
	user, err := userDB.GetUserByID(task.UserID)
	if err != nil {
		log.Fatalf("Error getting user by userID: %v", err)
		return err
	}

	// convert event timestamp to date string
	eventTime := time.Unix(p.EventTime, 0).Format("2006-01-02 15:04:05")

	log.Printf(" [*] Send Signal Email to User %s %s %s %s %s %s", user.Email, user.Name, task.Stock, eventTime, p.Signal, p.Strategy)
	notify.SendNewSignalEMail(user.Email, user.Name, task.Stock, eventTime, p.Signal, p.Strategy)
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
