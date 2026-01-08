package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"api/model"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	_ = godotenv.Load() 

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL n'est pas défini !")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB : %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Impossible de joindre la DB : %v", err)
	}

	log.Println("Connexion à PostgreSQL validée")

	_, err = model.New("One Piece", "Eiichiro Oda", 1997, "07/01/1997", "08/01/2026")
	if err != nil {
		log.Fatal(err)
	}

	_, err = model.New("One Piece", "Eiichiro Oda", 1997, "07/01/1997", "08/01/2026")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, world!")
		}),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Serveur HTTP démarré sur :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur : %v", err)
		}
	}()

	<-stop
	log.Println("Signal d'arrêt reçu, fermeture propre...")


	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erreur lors de l'arrêt du serveur HTTP : %v", err)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("Erreur lors de la fermeture de la DB : %v", err)
	}

	log.Println("Application arrêtée proprement")


	//const name = "Alice"

	// const tmpl = "Nom : {{ .Name }}. Auteur : {{ .Author }}. Annee : {{ .Years }}. Creation : {{ .Created_at }}. modification : {{ .Updated_at }}"

	// t, err := template.New("tmpl").Parse(tmpl)

	// if err != nil {

	// panic(err)

	// }

	// err = t.Execute(os.Stdout, tmpl)

	// if err != nil {
	// 	log.Println("c'est bon ")
	// panic(err)

	// }

}
