package service

import (
	"context"
	"errors"
	"fmt"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func getVectorDBConnUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("VECTORDB_USER"),
		os.Getenv("VECTORDB_PASSWORD"),
		os.Getenv("VECTORDB_HOST"),
		os.Getenv("VECTORDB_PORT"),
		os.Getenv("VECTORDB_DBNAME"),
	)
}
func openVectorStoreConn() (*pgvector.Store, error) {
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

func VectorizeCourseReviews(ctx context.Context, courseID int64) error {
	if courseID == 0 {
		return errors.New("course id is 0")
	}
	courseQuery := repository.NewCourseQuery()
	coursePO, err := courseQuery.GetCourse(ctx, courseQuery.WithID(courseID))
	if err != nil {
		fmt.Println(err)
		return err
	}

	courseName := coursePO.Name

	filter := domain.ReviewFilter{
		CourseID: courseID,
		Page:     0,
		PageSize: 100,
	}

	reviews, err := GetReviewList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var comments []string

	for _, review := range reviews {
		comments = append(comments, review.Comment)
	}

	vectorStore, err := openVectorStoreConn()

	if err != nil {
		fmt.Println(err)
		return err
	}

	targetStr := courseName + "\n" + strings.Join(comments, "\n")
	doc := schema.Document{
		PageContent: targetStr,
		Metadata: map[string]any{
			"courseID": courseID,
		},
	}

	_, err = vectorStore.AddDocuments(
		context.Background(),
		[]schema.Document{doc},
		vectorstores.WithReplacement(true),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = vectorStore.Close()
	return err
}

func GetMatchCourses(ctx context.Context, description string) ([]domain.Course, error) {
	vectorStore, err := openVectorStoreConn()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	docs, err := vectorStore.SimilaritySearch(context.Background(), description, 2)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = vectorStore.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var courseIDs []int64
	for _, doc := range docs {
		courseID := doc.Metadata["courseID"].(float64)
		courseIDs = append(courseIDs, int64(courseID))
	}

	query := repository.NewCourseQuery()

	coursePOs, err := query.GetCourseByIDs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	courseCategories, err := query.GetCourseCategories(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	reviewQuery := repository.NewReviewQuery()
	infos, err := reviewQuery.GetCourseReviewInfo(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]domain.Course, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCoursePOToDomain(coursePO)
		converter.PackCourseWithCategories(&course, courseCategories[int64(coursePO.ID)])
		converter.PackCourseWithReviewInfo(&course, infos[int64(coursePO.ID)])
		courses = append(courses, course)
	}
	return courses, nil

}
