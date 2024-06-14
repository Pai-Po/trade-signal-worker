package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID            string
	Name          string
	Email         string
	Password      string
	EmailVerified sql.NullString
	Image         sql.NullString
}

type UserDB interface {
	InsertUser(id string, name string, email string, password string, emailVerified string, image string) (User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	DeleteUserByID(id string) (int64, error) // returns number of rows deleted
	Close()
}

type PostgresUserDB struct {
	db *sql.DB
}

func NewPostgresUserDB(connStr string) (*PostgresUserDB, error) {
	if connStr == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Error loading .env file: %v", err)
		}
		connStr = os.Getenv("POSTGRES_URL")
		if connStr == "" {
			return nil, errors.New("POSTGRES_URL not set in .env file")
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	// Check if table exists, if not return error
	_, err = db.Exec(`SELECT 'public."User"'::regclass;`)
	if err != nil {
		return nil, fmt.Errorf("User table does not exist: %v", err)
	}

	return &PostgresUserDB{db: db}, nil
}
func (db *PostgresUserDB) DumpTableColumns(tableName string) error {
	rows, err := db.db.Query(`SELECT column_name FROM information_schema.columns WHERE table_name=$1`, tableName)
	if err != nil {
		log.Printf("Failed to query %s table: %v", tableName, err)
		return err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			log.Printf("Failed to scan column name: %v", err)
			return err
		}
		columns = append(columns, column)
	}

	for _, col := range columns {
		fmt.Println(col)
	}

	return nil
}
func (db *PostgresUserDB) GetUserByID(id string) (User, error) {
	// db.DumpTableColumns("User")
	row := db.db.QueryRow(`SELECT * FROM "User" WHERE "id" = $1;`, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.EmailVerified, &user.Image)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return User{}, err
	}

	return user, nil
}

func (db *PostgresUserDB) Close() {
	db.db.Close()
}
