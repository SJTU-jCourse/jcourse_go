package embedding

import (
	"context"
	"fmt"
	"log"
	"sync"

	"jcourse_go/util"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

var (
	vectorStoreInstance *pgvector.Store
	initOnce            sync.Once
	initErr             error
)

func InitVectorStore(ctx context.Context) error {
	initOnce.Do(func() {
		llm, err := openai.New(
			openai.WithBaseURL(util.GetOpenAIBaseURL()),
			openai.WithEmbeddingModel(util.GetEmbeddingModelName()),
			openai.WithToken(util.GetOpenAIToken()),
		)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize OpenAI client: %w", err)
			log.Printf("ERROR: %v. Check API key and connectivity.", initErr)
			return
		}

		embedder, err := embeddings.NewEmbedder(llm)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize OpenAI embedder: %w", err)
			log.Printf("ERROR: %v", initErr)
			return
		}
		log.Println("OpenAI embedder initialized.")

		connURL := util.GetVectorDBConnUrl()
		if connURL == "" {
			initErr = fmt.Errorf("vector DB connection URL is empty (set VECTOR_DB_URL)")
			log.Printf("ERROR: %v", initErr)
			return
		}
		log.Printf("Using Vector DB URL: %s", connURL)

		store, err := pgvector.New(
			ctx, // timeout 20s
			pgvector.WithConnectionURL(connURL),
			pgvector.WithEmbedder(embedder),
			pgvector.WithCollectionName("course_embeddings"),
		)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize pgvector store: %w", err)
			log.Printf("ERROR: %v. Ensure PostgreSQL is running, pgvector extension exists, and URL is correct.", initErr)
			return
		}

		log.Println("pgvector store initialized successfully.")
		vectorStoreInstance = &store
	})

	return initErr
}

func GetStore() (*pgvector.Store, error) {
	if initErr != nil || vectorStoreInstance == nil {
		return nil, fmt.Errorf("vectorbase is not available now")
	}
	return vectorStoreInstance, nil
}

func CloseVectorStore() error {
	log.Println("Closing vector store client resources (if any)...")
	vectorStoreInstance = nil
	initErr = nil
	log.Println("Vector store client resources handled.")
	return nil
}
