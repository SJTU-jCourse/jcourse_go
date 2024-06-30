package rpc

import (
	"github.com/sashabaranov/go-openai"
)

var llmClient *openai.Client

func InitOpenAIClient() {
	config := openai.ClientConfig{
		BaseURL:              "",
		OrgID:                "",
		APIType:              "",
		APIVersion:           "",
		AssistantVersion:     "",
		AzureModelMapperFunc: nil,
		HTTPClient:           nil,
		EmptyMessagesLimit:   0,
	}
	llmClient = openai.NewClientWithConfig(config)
}

func GetOpenAIClient() *openai.Client {
	return llmClient
}
