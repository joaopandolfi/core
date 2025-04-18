package core

import (
	"context"
)

// Embedder is an interface for generating vectors from content
type Embedder interface {
	// GenerateEmbedding generates vector embeddings based on input content
	GenerateEmbedding(ctx context.Context, content string) (*Embedding, error)
}
