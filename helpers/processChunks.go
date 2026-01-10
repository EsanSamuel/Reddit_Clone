package helpers

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/EsanSamuel/Reddit_Clone/config"
)

func ProcessChunks(allChunks []Chunk, queryEmbeddings []float32, query string) ([]float32, string) {
	for i := range allChunks {
		if allChunks[i].Embedding == nil {
			allChunks[i].Embedding = config.AIEmbeddings(allChunks[i].Text)
		}
	}

	//var scores []float32
	for i, chunk := range allChunks {
		score := CosineSimilarity(queryEmbeddings, chunk.Embedding)
		allChunks[i].Score = score
	}

	sort.Slice(allChunks, func(i, j int) bool {
		return allChunks[i].Score > allChunks[j].Score
	})

	var scores []float32
	for i := range allChunks {
		scores = append(scores, allChunks[i].Score)
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i] > scores[j]
	})

	topK := 3
	var content []string
	for i, chunk := range allChunks {
		if i < topK && chunk.Score > 0.35 {
			content = append(content, chunk.Text)
		}
	}

	contentString, err := json.Marshal(content)

	prompt := fmt.Sprintf(`You are an AI assistant. Use the following content to answer the user's query.

                    **Instructions:**
                        1. Only use the information provided in the relevant content chunks.
                        2. Provide a clear, concise, and informative answer.
                        3. Highlight disagreements, recurring ideas, or differing opinions if present.
                        4. Keep the tone neutral, factual, and professional.
                        5. Do not include information not present in the content.

                    **User Query: "%s" ** 
                    

                    **Relevant Content Chunks:"%s"** 
                        

                    **Answer:**
                     `, query, string(contentString))

	answer, err := config.Ai(prompt)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}
	fmt.Println(content)

	return scores, answer
}
