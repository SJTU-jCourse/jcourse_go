package rpc

import (
	"context"
	"fmt"
	"log"

	// "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
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

func HelloWorld() {
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, "The first man to walk on the moon",
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
