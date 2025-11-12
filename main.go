package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// Message represents a stored message
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageStore holds messages in memory
type MessageStore struct {
	messages []Message
	nextID   int
	mu       sync.RWMutex
}

var store = &MessageStore{
	messages: make([]Message, 0),
	nextID:   1,
}

// AddMessage adds a new message and returns it
func (ms *MessageStore) AddMessage(content string) Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	msg := Message{
		ID:        ms.nextID,
		Content:   content,
		Timestamp: time.Now(),
	}

	ms.messages = append(ms.messages, msg)
	ms.nextID++

	return msg
}

// GetRecentMessages returns the last N messages
func (ms *MessageStore) GetRecentMessages(count int) []Message {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	total := len(ms.messages)
	if total == 0 {
		return []Message{}
	}

	start := total - count
	if start < 0 {
		start = 0
	}

	// Return a copy to avoid race conditions
	result := make([]Message, total-start)
	copy(result, ms.messages[start:])

	// Reverse to show newest first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func main() {
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
			msg := store.AddMessage(req.Message)

			// Return success with the message
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": msg,
			})

		} else if r.Method == http.MethodGet {
			// Get last 3 messages
			messages := store.GetRecentMessages(3)

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
