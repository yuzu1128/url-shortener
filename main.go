package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// URLStore stores short URLs in memory
type URLStore struct {
	mu    sync.RWMutex
	urls  map[string]URLEntry
}

// URLEntry represents a shortened URL
type URLEntry struct {
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int       `json:"clicks"`
}

// Request/Response types
type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

type StatsResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var store = &URLStore{
	urls: make(map[string]URLEntry),
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/stats/", handleStats)

	log.Println("URL Shortener API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Handle redirect for short URLs
	if r.URL.Path != "/" {
		shortCode := r.URL.Path[1:]
		store.mu.RLock()
		entry, exists := store.urls[shortCode]
		store.mu.RUnlock()

		if !exists {
			http.NotFound(w, r)
			return
		}

		// Increment click count
		store.mu.Lock()
		entry.Clicks++
		store.urls[shortCode] = entry
		store.mu.Unlock()

		http.Redirect(w, r, entry.OriginalURL, http.StatusFound)
		return
	}

	// API info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service": "URL Shortener API",
		"version": "1.0.0",
		"endpoints": "POST /shorten, GET /{shortCode}, GET /stats/{shortCode}",
	})
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	if req.URL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "URL is required"})
		return
	}

	// Generate short code
	shortCode := generateShortCode()

	entry := URLEntry{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	store.mu.Lock()
	store.urls[shortCode] = entry
	store.mu.Unlock()

	// Get base URL from environment or default
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{
		ShortURL:    baseURL + "/" + shortCode,
		ShortCode:   shortCode,
		OriginalURL: req.URL,
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	shortCode := r.URL.Path[7:] // Remove "/stats/" prefix

	store.mu.RLock()
	entry, exists := store.urls[shortCode]
	store.mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Short URL not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatsResponse{
		ShortCode:   entry.ShortCode,
		OriginalURL: entry.OriginalURL,
		Clicks:      entry.Clicks,
		CreatedAt:   entry.CreatedAt,
	})
}

func generateShortCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
