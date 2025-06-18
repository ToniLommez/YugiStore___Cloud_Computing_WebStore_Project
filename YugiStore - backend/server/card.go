package server

import (
	"context"
	"log"

	"google.golang.org/api/iterator"
)

type Card struct {
	Name     string `json:"name" firestore:"name"`
	Edition  string `json:"edition" firestore:"edition"`
	ImageURL string `json:"image_url" firestore:"image_url"`
}

func CreateCard(ctx context.Context, card Card) error {
	_, err := FirestoreClient.Collection("cards").Doc(card.Edition).Set(ctx, card)
	if err != nil {
		log.Printf("Failed to save card with ID '%s': %v", card.Edition, err)
		return err
	}

	return nil
}

func GetCardsByName(ctx context.Context, name string) ([]Card, error) {
	var cards []Card

	iter := FirestoreClient.Collection("cards").Where("name", "==", name).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var card Card
		if err := doc.DataTo(&card); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func GetAllCards(ctx context.Context) ([]Card, error) {
	var cards []Card

	iter := FirestoreClient.Collection("cards").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var card Card
		if err := doc.DataTo(&card); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func GetAllEditionsName(ctx context.Context) ([]string, error) {
	var editions []string

	iter := FirestoreClient.Collection("cards").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		edition, err := doc.DataAt("edition")
		name, err := doc.DataAt("name")
		if err != nil {
			continue
		}

		if editionStr, ok := edition.(string); ok {
			if nameStr, ok := name.(string); ok {
				editions = append(editions, nameStr+" - "+editionStr)
			}
		}
	}

	return editions, nil
}
