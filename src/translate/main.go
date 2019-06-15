// Sample translate-quickstart translates "Hello, world!" into Russian.
package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the text to translate.
	text := "ఈరోజు ఎన్టీఆర్ జయంతి సందర్భంగా ఆ"
	// Sets the target language.
	target, err := language.Parse("en")
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	// Translates the text into Telugu.
	translations, err := client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	detect, err := detectLanguage(text)
	if err != nil {
		log.Fatal("Failed to detect the language: %v", err)
	}
	fmt.Printf("Detected the language: %v\n", detect)

	fmt.Printf("Text: %v\n", text)
	fmt.Printf("Translation: %v\n", translations[0].Text)
}

func detectLanguage(text string) (*translate.Detection, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	lang, err := client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return &lang[0][0], nil
}
