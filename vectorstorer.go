package core

import "context"

// Vec32 represents a generic vector of any dimension
// that is quantized as a 32-bit floating point
type Vec32 []float32

// Embedding pairs a vector with an identifier and its content
type Embedding struct {
	ID      string
	Vector  Vec32
	Content string
}

// SearchResult represents a single result from a vector search
type SearchResult struct {
	Score      float32
	Embedding  *Embedding
	SearchMeta *SearchParams
}

// SearchParams contains parameters for vector search operations
type SearchParams struct {
	Query     string
	QueryVec  Vec32
	Limit     int
	Threshold float32
}

type VectorStorer interface {
	// Add stores embeddings in the database
	Add(ctx context.Context, contents []string) ([]*Embedding, error)

	// Search finds vectors similar to the query vector
	Search(ctx context.Context, params *SearchParams) ([]*SearchResult, error)

	// Close releases resources associated with the vector storer
	Close() error
}
