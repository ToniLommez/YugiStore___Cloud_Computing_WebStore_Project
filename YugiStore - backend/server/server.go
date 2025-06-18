package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	FirestoreClient *firestore.Client
	StorageClient   *storage.Client
)

func setupRoutes() {
	http.HandleFunc("/ping", withCORS(ping))
	http.HandleFunc("/card", withCORS(cardsHandler))
	http.HandleFunc("/editions", withCORS(editionsHandler))
	http.HandleFunc("/product", withCORS(productsHandler))
}

func StartServer() {
	log.Print("starting server...")

	ctx := context.Background()
	sa := option.WithCredentialsFile("yugistore-1c7f8e32ee8e.json")

	// Firebase App
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	// Firestore
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer fsClient.Close()
	FirestoreClient = fsClient
	log.Printf("Connected to Firestore")

	// Cloud Storage
	stClient, err := storage.NewClient(ctx, sa)
	if err != nil {
		log.Fatalln(err)
	}
	defer stClient.Close()
	StorageClient = stClient
	log.Printf("Connected to Cloud Storage")

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start server
	setupRoutes()
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
