package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"echo/internal"
	"echo/telegram"
)

type TelegramButton struct {
	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}
type MessageRequest struct {
	TopicName *string         `json:"topicName,omitempty"`
	Title     *string         `json:"title,omitempty"`
	Message   string          `json:"message"`
	Enqueue   bool            `json:"enqueue"`
	Time      *int64          `json:"time,omitempty"`
	Button    *TelegramButton `json:"button,omitempty"`
}

type JsonResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) HandleMessage(w http.ResponseWriter, r *http.Request) {
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	message := req.Message
	uniqueId := internal.GenerateMD5Hash(message)

	ctx := context.Background()
	rd := app.config.redis

	exist, err := rd.KeyExists(ctx, uniqueId)
	if err != nil {
		log.Printf("Error checking key existence: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if exist {
		log.Println("Message already exists in Redis")
		respondWithError(w, "Message already exists", http.StatusConflict)
		return
	}

	if req.Enqueue {
		ttl := time.Duration(*req.Time) * time.Second
		err = rd.Enqueue(ctx, uniqueId, message, ttl)
		if err != nil {
			log.Printf("Could not enqueue message: %v", err)
			http.Error(w, "Could not enqueue message", http.StatusInternalServerError)
			return
		}

	}

	telegram, err := telegram.NewTelegramBot()
	if err != nil {
		log.Printf("Could not initialize the Telegram client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var buttonText *string
	var buttonURL *string

	// Only set buttonText and buttonURL if button is provided in the request
	if req.Button != nil {
		buttonText = &req.Button.Text
		buttonURL = &req.Button.URL
	}
	if err := telegram.SendMessage(req.Message, buttonText, buttonURL); err != nil {
		log.Printf("Failed to send message to Telegram: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, JsonResponse{Message: "Message processed", Status: "success"}, http.StatusOK)
}

func (app *application) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}
func (app *application) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	rd := app.config.redis
	ctx := context.Background()
	messages, err := rd.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get all the messages: %v", err)
		http.Error(w, "Failed to retrieve all the messages", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, JsonResponse{Data: json.RawMessage(messages), Status: "success"}, http.StatusOK)
}
func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	respondWithError(w, "Page not found", http.StatusNotFound)
}
func respondWithError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JsonResponse{Message: message, Status: "error"})
}

func respondWithJSON(w http.ResponseWriter, payload JsonResponse, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
