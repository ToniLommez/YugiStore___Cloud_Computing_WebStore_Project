package server

import (
	"context"
	"log"

	"google.golang.org/api/iterator"
)

type Product struct {
	Edition string `json:"edition" firestore:"edition"`
	State   string `json:"state" firestore:"state"`
	Price   string `json:"price" firestore:"price"`
}

func RegisterProduct(ctx context.Context, product Product) error {
	_, _, err := FirestoreClient.Collection("products").Add(ctx, product)
	if err != nil {
		log.Printf("Failed to add product: %v", err)
		return err
	}

	return nil
}

func GetProducts(ctx context.Context, edition string) ([]Product, error) {
	var products []Product

	iter := FirestoreClient.Collection("products").Where("edition", "==", edition).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var product Product
		if err := doc.DataTo(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetAllProducts(ctx context.Context) ([]Product, error) {
	var products []Product

	iter := FirestoreClient.Collection("products").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var product Product
		if err := doc.DataTo(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
