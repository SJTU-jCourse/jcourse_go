package embedding

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai" // Using OpenAI
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

var (
	vectorStoreInstance *pgvector.Store
	initOnce            sync.Once
	initErr             error
)

func InitVectorStore(ctx context.Context) error {
	initOnce.Do(func() {
		log.Println("Initializing singleton vector store client with OpenAI...")

		llm, err := openai.New() // TODO set url model and token
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

		connURL := getVectorDBConnUrl()
		if connURL == "" {
			initErr = fmt.Errorf("vector DB connection URL is empty (set VECTOR_DB_URL)")
			log.Printf("ERROR: %v", initErr)
			return
		}
		log.Printf("Using Vector DB URL: %s", connURL)

		store, err := pgvector.New(
			context.Background(), // Store lifecycle tied to application
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

func GetStore() *pgvector.Store {
	return vectorStoreInstance
}

func CloseVectorStore() error {
	log.Println("Closing vector store client resources (if any)...")
	vectorStoreInstance = nil
	initErr = nil
	log.Println("Vector store client resources handled.")
	return nil
}

func getVectorDBConnUrl() string {
	dbURL := os.Getenv("VECTOR_DB_URL") // TODO set vector db url
	if dbURL == "" {
		dbURL = "postgresql://user:password@localhost:5432/vector_db?sslmode=disable" // !! REPLACE DEFAULT !!
		log.Printf("WARNING: VECTOR_DB_URL env var not set. Using default: %s", dbURL)
	}
	return dbURL
}
