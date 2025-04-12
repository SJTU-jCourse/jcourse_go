package embedding

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

const pgvectorURL = "postgres://test_user:your_secure_password@localhost:5432/test_vector_db?sslmode=disable"

func TestPgvectorStoreRest(t *testing.T) {
	ctx := context.Background()

	llm, err := openai.New(
		openai.WithBaseURL("https://api.oaipro.com/v1"),
		openai.WithToken("sk-awcK02dXpXqQ4Kew86D24305Ff824f29A0162f4142F89d4f"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	)
	require.NoError(t, err)
	e, err := embeddings.NewEmbedder(llm)
	require.NoError(t, err)

	conn, err := pgx.Connect(ctx, pgvectorURL)
	require.NoError(t, err)

	store, err := pgvector.New(
		ctx,
		pgvector.WithConn(conn),
		pgvector.WithEmbedder(e),
		pgvector.WithPreDeleteCollection(true),
		pgvector.WithEmbeddingTableName("test_embedding"),
		//pgvector.WithCollectionName(makeNewCollectionName()),
	)
	require.NoError(t, err)

	_, err = store.AddDocuments(ctx, []schema.Document{
		{PageContent: "大学英语 I", Metadata: map[string]any{
			"country": "japan",
		}},
		{PageContent: "大学英语 II"},
		{PageContent: "高等数学 I"},
		{PageContent: "线性代数 II"},
	})
	require.NoError(t, err)

	docs, err := store.SimilaritySearch(ctx, "数学课", 2)
	require.NoError(t, err)
	require.Len(t, docs, 2)
	t.Logf("%+v", docs)
}
