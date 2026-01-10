package helpers

import "strings"

type Chunk struct {
	PostID    string
	ChunkId   int
	Source    string
	Text      string
	Embedding []float32
	Score     float32
}

func ChunkText(postId, source, text string) []Chunk {
	paragraphs := strings.Split(text, "\n")

	var Chunks []Chunk

	chunkId := 0

	for _, p := range paragraphs {
		if strings.TrimSpace(p) == "" {
			continue
		}
		Chunks = append(Chunks, Chunk{
			PostID:  postId,
			ChunkId: chunkId + 1,
			Source:  source,
			Text:    p,
		})
		chunkId++
	}
	return Chunks
}
