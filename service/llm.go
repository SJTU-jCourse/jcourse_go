package service

import (
	"context"
	"errors"
	"fmt"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"
)

func VectorizeCourseReviews(ctx context.Context, courseID int64) ([]float32, error) {
	if courseID == 0 {
		return nil, errors.New("course id is 0")
	}
	courseQuery := repository.NewCourseQuery()
	coursePO, err := courseQuery.GetCourse(ctx, courseQuery.WithID(courseID))
	if err != nil {
		return nil, err
	}

	courseName := coursePO.Name

	filter := domain.ReviewFilter{
		CourseID: courseID,
		Page:     0,
		PageSize: 100,
	}

	reviews, err := GetReviewList(ctx, filter)
	if err != nil {
		return nil, err
	}

	var comments []string

	for _, review := range reviews {
		comments = append(comments, review.Comment)
	}

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

	targetStr := courseName + "\n" + strings.Join(comments, "\n")
	embs, err := embedder.EmbedQuery(context.Background(), targetStr)

	fmt.Println(targetStr)
	// embs的维度
	fmt.Println(len(embs))

	// 将向量存入数据库

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	query := repository.NewCourseVectorQuery()

	err = query.UpdateCourseVector(ctx, courseID, embs)
	if err == gorm.ErrRecordNotFound {
		err = query.InsertCourseVector(ctx, courseID, embs)
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return embs, nil
}
