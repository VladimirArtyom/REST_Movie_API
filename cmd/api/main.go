package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dsn string
	}
}

type application struct {
	config config  
	logger *log.Logger 
} 

func main() {

	// Define the config object
	// Define flags for the config values
	// Parse the flags
	// Define Log object
	// Integrate it to application

	var cfg config
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.IntVar(&cfg.port, "port", 8080, "API Default server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (devevelopment|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	var logger *log.Logger = log.New(os.Stdout, "",  log.Ldate | log.Ltime)

	var app *application = &application{
		config: cfg,
		logger: logger,
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	print("database connection established\n")

	// Make a multiplexer, basically request handler using mux basically

	var server *http.Server = &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting server on %s in %s mode", server.Addr, cfg.env)
	server.ListenAndServe()

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create 5-second context
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second) 
	defer cancelFunc()
	
	// Check or ping the database availaility
	err = db.PingContext(ctx)
	if err != nil {
	return nil, err
	}

	return db, nil
}
