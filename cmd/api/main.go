package main

import (
	"context"
	"net/http"

	"github.com/MrTomatePNG/projeto-m/internal/database"
	"github.com/MrTomatePNG/projeto-m/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgxpool.New(context.Background(), "postgresql://postgres:postgres123@localhost:5432/meuapp_dev")
	if err != nil {
		panic(err)
	}
	queries := database.New(conn)

	handlers := handlers.NewUserHandler(queries)

	r := chi.NewRouter()

	r.Post("/register", handlers.Create())
	r.Get("/login", handlers.Login())
	http.ListenAndServe(":8080", r)
	defer conn.Close()
}
