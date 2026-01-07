package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"context"
	"time"

	  "github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// Charge le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Pas de fichier .env trouvé, les variables d'environnement doivent être définies autrement")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL n'est pas défini !")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB : %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

		if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Impossible de joindre la DB : %v", err)
	}

	fmt.Println("Connexion réussie à PostgreSQL !")

	var greeting string
	err = db.QueryRow("SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}

	fmt.Println(greeting)
}
