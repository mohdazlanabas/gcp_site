package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"example.com/my-app/database"
	"github.com/jackc/pgx/v5"
)

// Message represents a stored message
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// dbConn is the global database connection pool
var dbConn *pgx.Conn

// AddMessage adds a new message to the database
func AddMessage(ctx context.Context, content string) (Message, error) {
	msg := Message{
		Content:   content,
		Timestamp: time.Now(),
	}

	query := "INSERT INTO messages (content, timestamp) VALUES ($1, $2) RETURNING id"
	err := dbConn.QueryRow(ctx, query, msg.Content, msg.Timestamp).Scan(&msg.ID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to insert message: %w", err)
	}

	return msg, nil
}

// GetRecentMessages returns the last N messages from the database
func GetRecentMessages(ctx context.Context, count int) ([]Message, error) {
	var messages []Message
	query := "SELECT id, content, timestamp FROM messages ORDER BY timestamp DESC LIMIT $1"
	rows, err := dbConn.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through message rows: %w", err)
	}

	// Reverse the slice to show oldest first (as per original in-memory store behavior)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func main() {
	// Initialize database connection
	var err error
	dbConn, err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if dbConn != nil {
			dbConn.Close(context.Background())
			log.Println("Database connection closed.")
		}
	}()

	// Run database migrations
	err = database.MigrateDB(dbConn)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Serve static files from the frontend folder under /static/
	fs := http.FileServer(http.Dir(filepath.Join("src", "frontend")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve root index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("src", "frontend", "index.html"))
	})

	// Simple API endpoint
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{"message": "\n\nHello from the Go backend on App Engine"}
		_ = json.NewEncoder(w).Encode(resp)
	})

	// POST /api/messages - Save a new message
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			// Parse request body
			var req struct {
				Message string `json:"message"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
				return
			}

			if req.Message == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Message cannot be empty"})
				return
			}

			// Add message to store
			msg, err := AddMessage(r.Context(), req.Message)
			if err != nil {
				log.Printf("Error adding message to DB: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save message"})
				return
			}

			// Return success with the message
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": msg,
			})

		} else if r.Method == http.MethodGet {
			// Get last 3 messages
			messages, err := GetRecentMessages(r.Context(), 3)
			if err != nil {
				log.Printf("Error getting recent messages from DB: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve messages"})
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"messages": messages,
				"count":    len(messages),
			})

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	})

	addr := ":8080" // App Engine will override the port
	log.Printf("Server starting on %s", addr)
	log.Printf("Message store initialized")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
