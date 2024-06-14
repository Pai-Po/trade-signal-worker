package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Task struct {
	ID           int
	UserID       string
	Stock        string
	KlineType    string
	BuyStrategy  string
	SellStrategy string
	Status       string
	Timestamp    string
}

type TaskDB interface {
	InsertTask(userID string, stock string, klineType string, buyStrategy string, sellStrategy string, status string) (Task, error)
	GetAllTasks() ([]Task, error)
	GetTaskByID(userID string, taskID int) (Task, error)
	GetTaskByTaskID(taskID int) (Task, error)
	GetTasksByUserID(userID string) ([]Task, error)
	GetAllRunningTasks() ([]Task, error)
	DeleteAllTasks() (int64, error)                           // returns number of rows deleted
	DeleteTasksByID(userID string, taskID int) (int64, error) // returns number of rows deleted
	DeleteTable() error
	Close()
}

type PostgresTaskDB struct {
	db *sql.DB
}

func NewPostgresTaskDB(connStr string) *PostgresTaskDB {
	if connStr == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
			return nil
		}
		connStr = os.Getenv("POSTGRES_URL")
		if connStr == "" {
			log.Fatalf("POSTGRES_URL not set in .env file")
			return nil
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id serial PRIMARY KEY, 
            user_id varchar,
            stock varchar,
            kline_type varchar,
            buy_strategy varchar,
            sell_strategy varchar,
            status varchar,
            timestamp timestamp default current_timestamp
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create task table: %v", err)
	}

	return &PostgresTaskDB{db: db}
}

func (db *PostgresTaskDB) InsertTask(userID string, stock string, klineType string, buyStrategy string, sellStrategy string, status string) (Task, error) {
	_, err := db.db.Exec(`
        INSERT INTO tasks (user_id, stock, kline_type, buy_strategy, sell_strategy, status) 
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, userID, stock, klineType, buyStrategy, sellStrategy, status)
	if err != nil {
		log.Printf("Failed to insert task: %v", err)
		return Task{}, err
	}

	var task Task
	err = db.db.QueryRow("SELECT * FROM tasks WHERE id = (SELECT MAX(id) FROM tasks);").Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
	if err != nil {
		log.Printf("Failed to fetch inserted task: %v", err)
		return Task{}, err
	}

	return task, nil
}

func (db *PostgresTaskDB) GetAllTasks() ([]Task, error) {
	rows, err := db.db.Query("SELECT * FROM tasks;")
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
		if err != nil {
			log.Printf("Failed to scan task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *PostgresTaskDB) GetTaskByID(userID string, taskID int) (Task, error) {
	row := db.db.QueryRow("SELECT * FROM tasks WHERE user_id = $1 AND id = $2;", userID, taskID)

	var task Task
	err := row.Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
	if err != nil {
		log.Printf("Failed to fetch task: %v", err)
		return Task{}, err
	}

	return task, nil
}

func (db *PostgresTaskDB) GetTaskByTaskID(taskID int) (Task, error) {
	row := db.db.QueryRow("SELECT * FROM tasks WHERE id = $1;", taskID)

	var task Task
	err := row.Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
	if err != nil {
		log.Printf("Failed to fetch task: %v", err)
		return Task{}, err
	}

	return task, nil
}

func (db *PostgresTaskDB) GetTasksByUserID(userID string) ([]Task, error) {
	rows, err := db.db.Query("SELECT * FROM tasks WHERE user_id = $1;", userID)
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
		if err != nil {
			log.Printf("Failed to scan task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *PostgresTaskDB) GetAllRunningTasks() ([]Task, error) {
	rows, err := db.db.Query("SELECT * FROM tasks WHERE status = 'running';")
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Stock, &task.KlineType, &task.BuyStrategy, &task.SellStrategy, &task.Status, &task.Timestamp)
		if err != nil {
			log.Printf("Failed to scan task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (db *PostgresTaskDB) DeleteAllTasks() (int64, error) {
	res, err := db.db.Exec("DELETE FROM tasks;")
	if err != nil {
		log.Printf("Failed to delete tasks: %v", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v", err)
		return 0, err
	}

	return rowsAffected, nil
}

func (db *PostgresTaskDB) DeleteTaskByID(userID string, taskID int) (int64, error) {
	res, err := db.db.Exec("DELETE FROM tasks WHERE user_id = $1 AND id = $2;", userID, taskID)
	if err != nil {
		log.Printf("Failed to delete task: %v", err)
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v", err)
		return 0, err
	}

	return rowsAffected, nil
}

func (db *PostgresTaskDB) DeleteTable() error {
	_, err := db.db.Exec("DROP TABLE IF EXISTS tasks;")
	if err != nil {
		log.Printf("Failed to delete table: %v", err)
		return err
	}

	return nil
}

func (db *PostgresTaskDB) Close() {
	db.db.Close()
}
