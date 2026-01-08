package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"api/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL n'est pas défini !")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("DB unreachable: %v", err)
	}

	log.Println("Connexion PostgreSQL OK")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
		SELECT id, title, author, years, created_at, updated_at
		FROM books`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		books := []model.Book{}

		for rows.Next() {
			var b model.Book
			err := rows.Scan(
				&b.Id,
				&b.Title,
				&b.Author,
				&b.Years,
				&b.CreatedAt,
				&b.UpdatedAt,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			books = append(books, b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	})

	r.Post("/books", func(w http.ResponseWriter, r *http.Request) {
		book := &model.Book{}

		createdBook, err := model.CreateBookFromRequest(book, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBook)
	})

	r.Put("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		book := &model.Book{
			Id: id,
		}

		if err := model.UpdateBookFromRequest(book, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Serveur HTTP démarré sur :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Serveur error: %v", err)
		}
	}()

	<-stop
	log.Println("Arrêt du serveur...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	server.Shutdown(shutdownCtx)

	log.Println("Application arrêtée proprement")

	
}
