package rpc

import (
	"context"
	"fmt"

	"jcourse_go/util"

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
		util.GetPostgresUser(),
		util.GetPostgresPassword(),
		util.GetPostgresHost(),
		util.GetPostgresPort(),
		util.GetPostgresPassword(),
	)
}
func OpenVectorStoreConn(ctx context.Context) (*pgvector.Store, error) {
	llm, err := openai.New()
	if err != nil {

		return nil, err
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {

		return nil, err
	}

	store, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(getVectorDBConnUrl()),
		pgvector.WithEmbedder(embedder),
	)

	return &store, err
}
