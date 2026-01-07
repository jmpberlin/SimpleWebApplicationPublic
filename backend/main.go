package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getPort() string {
	return getEnv("PORT", "8081")
}

func initDB() error {
	host := getEnv("POSTGRES_HOST", "postgres")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "")
	dbname := getEnv("POSTGRES_DB", "myapp")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("Connecting to database at %s:%s/%s...", host, port, dbname)

	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✓ Database connection established!")
				return nil
			}
		}
		log.Printf("Database connection attempt %d/10 failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after 10 attempts: %v", err)
}

func getMessageByID(id int) (string, error) {
	var message string
	err := db.QueryRow("SELECT message FROM messages WHERE id = $1", id).Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s - Remote: %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

	message, err := getMessageByID(1)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", message)
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s - Remote: %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

	message, err := getMessageByID(2)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", message)
}

func langingPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s - Remote: %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

	message, err := getMessageByID(3)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", message)
}

func impressumHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s %s - Remote: %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
	fmt.Fprintf(w, "Copyright (c) 2025 by github.com/jmpberlin")
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	port := getPort()

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/bye", byeHandler)
	http.HandleFunc("/impressum", impressumHandler)
	http.HandleFunc("/", langingPageHandler)

	log.Printf("Server starting on port %s...", port)
	log.Println("Available routes:")
	log.Printf("  - http://localhost:%s/hello", port)
	log.Printf("  - http://localhost:%s/bye", port)
	log.Printf("  - http://localhost:%s/", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
