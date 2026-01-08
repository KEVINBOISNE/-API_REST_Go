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
	"api/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// Chargement des variables d'environnement
	_ = godotenv.Load()

	// Connexion à la DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL n'est pas défini")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Erreur DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Erreur fermeture DB: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("DB inaccessible: %v", err)
	}

	log.Println("Connexion PostgreSQL OK")

	// Repository
	repo := repository.NewBookRepositorySQL(db)

	// Router Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	/* ---------- GET ALL ---------- */
	r.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := repo.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	})

	/* ---------- GET BY ID ---------- */
	r.Get("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		book, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Livre non trouvé", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})

	/* ---------- POST ---------- */
	r.Post("/books", func(w http.ResponseWriter, r *http.Request) {
    // PAS besoin de créer un Book vide avant
    book, err := model.CreateBookFromRequest(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Ajoute les timestamps
    book.CreatedAt = time.Now().Format("2006-01-02")
    book.UpdatedAt = book.CreatedAt

    if err := repo.Create(r.Context(), book); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
})


	/* ---------- PUT ---------- */
	r.Put("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		book, err := repo.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Livre non trouvé", http.StatusNotFound)
			return
		}

		if err := model.UpdateBookFromRequest(book, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		book.UpdatedAt = time.Now().Format("2006-01-02")

		if err := repo.Update(r.Context(), book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	})

	/* ---------- DELETE ---------- */
	r.Delete("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := repo.Delete(r.Context(), id); err != nil {
			http.Error(w, "Livre non trouvé", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// Serveur HTTP
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Lancement serveur dans goroutine
	go func() {
		log.Println("Serveur HTTP démarré sur :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	// Gestion du signal d'arrêt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Arrêt du serveur...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erreur arrêt serveur: %v", err)
	}

	log.Println("Application arrêtée proprement")
}
