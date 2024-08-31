package rpc

import (
	"context"
	"fmt"
	"os"

	// "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

// var llmClient *openai.Client

// func InitOpenAIClient() {
// 	config := openai.ClientConfig{
// 		BaseURL:              "",
// 		OrgID:                "",
// 		APIType:              "",
// 		APIVersion:           "",
// 		AssistantVersion:     "",
// 		AzureModelMapperFunc: nil,
// 		HTTPClient:           nil,
// 		EmptyMessagesLimit:   0,
// 	}
// 	llmClient = openai.NewClientWithConfig(config)
// }

// func GetOpenAIClient() *openai.Client {
// 	return llmClient
// }

func getVectorDBConnUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("VECTORDB_USER"),
		os.Getenv("VECTORDB_PASSWORD"),
		os.Getenv("VECTORDB_HOST"),
		os.Getenv("VECTORDB_PORT"),
		os.Getenv("VECTORDB_DBNAME"),
	)
}
func OpenVectorStoreConn() (*pgvector.Store, error) {
	llm, err := openai.New()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	store, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(getVectorDBConnUrl()),
		pgvector.WithEmbedder(embedder),
	)

	return &store, err
}
