package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func cardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if name != "" {
			handleGetCardByName(w, r)
		} else {
			handleGetAllCards(w, r)
		}
	case http.MethodPost:
		handleAddCard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func editionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetEditionsName(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// Convert JSON
	name := r.FormValue("name")
	edition := r.FormValue("edition")
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload image to Cloud Storage
	filename := fmt.Sprintf("%s-%s", name, edition)
	imageURL, err := UploadImage(ctx, file, filename)
	if err != nil {
		log.Printf("Image upload failed: %v", err)
		http.Error(w, "Image upload failed", http.StatusInternalServerError)
		return
	}

	// Create Card to firestore
	card := Card{
		Name:     name,
		Edition:  edition,
		ImageURL: imageURL,
	}

	err = CreateCard(ctx, card)
	if err != nil {
		log.Printf("Failed to create card: %v", err)
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	// Write Status
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "card saved",
	})
}

func handleGetCardByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := r.URL.Query().Get("name")

	cards, err := GetCardsByName(ctx, name)
	if err != nil {
		log.Printf("Failed to get cards by name: %v", err)
		http.Error(w, "Failed to retrieve cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func handleGetAllCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cards, err := GetAllCards(ctx)
	if err != nil {
		log.Printf("Failed to get all cards: %v", err)
		http.Error(w, "Failed to retrieve cards", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func handleGetEditionsName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	editions, err := GetAllEditionsName(ctx)
	if err != nil {
		log.Printf("Failed to get all editions: %v", err)
		http.Error(w, "Failed to retrieve editions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editions)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		edition := r.URL.Query().Get("edition")
		if edition != "" {
			handleGetProduct(w, r)
		} else {
			handleGetAllProducts(w, r)
		}
	case http.MethodPost:
		handleAddProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Register product on firestore
	err := RegisterProduct(ctx, product)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Write Status
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "product saved",
	})
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	edition := r.URL.Query().Get("edition")

	products, err := GetProducts(ctx, edition)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := GetAllProducts(ctx)
	if err != nil {
		log.Printf("Failed to get all products: %v", err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
