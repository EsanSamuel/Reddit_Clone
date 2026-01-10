package config

import (
	"context"

	"log"
	"os"

	"google.golang.org/genai"
)

func Ai(prompt string) (string, error) {
	ctx := context.Background()

	api_key := os.Getenv("GEMINI_API_KEY")
	config := &genai.ClientConfig{
		APIKey: api_key,
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	//fmt.Println(result.Text())
	return result.Text(), nil
}

func AIEmbeddings(content string) []float32 {
	api_key := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()
	config := &genai.ClientConfig{
		APIKey: api_key,
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText(content, genai.RoleUser),
	}
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	//embeddings, err := json.MarshalIndent(result.Embeddings, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(embeddings))

	return result.Embeddings[0].Values
}
