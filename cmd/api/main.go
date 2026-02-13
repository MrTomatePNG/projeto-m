package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/MrTomatePNG/projeto-m/internal/auth"
	"github.com/MrTomatePNG/projeto-m/internal/database"
	"github.com/MrTomatePNG/projeto-m/internal/handlers"
	"github.com/MrTomatePNG/projeto-m/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn := initDB()
	defer conn.Close()

	queries := database.New(conn)

	jwtm, err := auth.NewJWTManager(os.Getenv("JWT_SECRET"), 24*time.Hour)
	if err != nil {
		panic(err)
	}

	userHandler := handlers.NewUserHandler(queries, jwtm)

	r := chi.NewRouter()

	r.Post("/register", userHandler.Create())
	r.Post("/login", userHandler.Login())

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtm))
		r.Get("/me", userHandler.Me())
	})
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
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
