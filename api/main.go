package main

import (
	"echo/internal"
	"echo/redis"

	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}
type config struct {
	addr  string
	redis *redis.RedisClient
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.HandleHealth)
	mux.HandleFunc("GET /message", app.GetAllMessages)
	mux.HandleFunc("POST /message", app.HandleMessage)

	// ignore not defined routes
	mux.HandleFunc("/", app.NotFound)
	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}
	log.Printf("Server starting in port: %s", app.config.addr)
	return server.ListenAndServe()

}
func main() {
	portApp, err := internal.GetEnv("PORT")
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalf("Could not initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	cfg := config{
		addr:  ":" + portApp,
		redis: redisClient,
	}

	app := &application{
		config: cfg,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
