package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
)

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

	addr := ":8080" // App Engine will override the port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
