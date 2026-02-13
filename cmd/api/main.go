package main

import (
	"context"
	"net/http"
	"os"

	"github.com/MrTomatePNG/projeto-m/internal/database"
	"github.com/MrTomatePNG/projeto-m/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn := initDB()

	queries := database.New(conn)

	handlers := handlers.NewUserHandler(queries)

	r := chi.NewRouter()

	r.Post("/register", handlers.Create())
	r.Get("/login", handlers.Login())
	http.ListenAndServe(":8080", r)
	defer conn.Close()
}

func initDB() *pgxpool.Pool {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		panic("DATABASE_URL must be set")
	}
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic("cannot conect database: " + err.Error())
	}
	return conn
}
