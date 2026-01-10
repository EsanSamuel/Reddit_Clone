package helpers

import (
	"fmt"
	"math"
)

func CosineSimilarity(vectorA []float32, vectorB []float32) float32 {
	if len(vectorA) != len(vectorB) {
		fmt.Println("Vector A and Vector B should be the same length")
		return 0
	}

	var dotProduct float32 = 0
	var magnitudeA float32 = 0
	var magnitudeB float32 = 0

	for i := 0; i < len(vectorA); i++ {
		dotProduct += vectorA[i] * vectorB[i]
		magnitudeA += vectorA[i] * vectorA[i]
		magnitudeB += vectorB[i] * vectorB[i]
	}

	magnitudeA = float32(math.Sqrt(float64(magnitudeA)))
	magnitudeB = float32(math.Sqrt(float64(magnitudeB)))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
