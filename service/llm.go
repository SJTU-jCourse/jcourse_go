package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/repository"
	"jcourse_go/rpc"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func OptCourseReview(courseName string, reviewContent string) (dto.OptCourseReviewResponse, error) {
	llm, err := openai.New()
	if err != nil {
		fmt.Println(err)
		return dto.OptCourseReviewResponse{}, err
	}
	// 构造提示词生成对话
	inputJson, _ := json.Marshal(map[string]string{
		"course": courseName,
		"review": reviewContent,
	})

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, constant.OptCourseReviewPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, string(inputJson)),
	}

	completion, err := llm.GenerateContent(
		context.Background(),
		content,
	)

	if err != nil {
		fmt.Println(err)
		return dto.OptCourseReviewResponse{}, err
	}
	var response dto.OptCourseReviewResponse
	err = json.Unmarshal([]byte(completion.Choices[0].Content), &response)

	return response, err
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

	vectorStore, err := rpc.OpenVectorStoreConn()

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
	vectorStore, err := rpc.OpenVectorStoreConn()

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
