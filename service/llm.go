package service

import (
	"context"
	"encoding/json"
	"fmt"
	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"jcourse_go/rpc"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// OptCourseReview使用LLM提示词对课程评价内容进行优化。
// courseName为课程名称，
// reviewContent为评价的内容，
// 函数返回值包含两个字段：
// Suggestion为修改建议，
// Result为根据修改建议给出的一种修改结果。
func OptCourseReview(courseName string, reviewContent string) (dto.OptCourseReviewResponse, error) {
	llm, err := openai.New()
	if err != nil {
		fmt.Println(err)
		return dto.OptCourseReviewResponse{}, err
	}
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

// GetCourseSummary使用LLM提示词基于课程评价生成课程总结。
// courseID为课程的ID，
// 返回值包含课程的总结。
// TODO: 此处基于课程最近的100条评价内容生成课程总结，后续可
// 以进一步调整和优化。
func GetCourseSummary(ctx context.Context, courseID int64) (*dto.GetCourseSummaryResponse, error) {
	courseQuery := repository.NewCourseQuery(dal.GetDBClient())
	coursePOs, err := courseQuery.GetCourse(ctx, repository.WithID(courseID))
	if err != nil || len(coursePOs) == 0 {
		return nil, err
	}
	coursePO := coursePOs[0]

	offeredCourseQuery := repository.NewOfferedCourseQuery(dal.GetDBClient())
	offeredCoursePOs, err := offeredCourseQuery.GetOfferedCourseTeacherGroup(ctx, []int64{courseID})
	if err != nil {
		return nil, err
	}

	reviewQuery := repository.NewReviewQuery(dal.GetDBClient())
	infos, err := reviewQuery.GetCourseReviewInfo(ctx, []int64{courseID})
	if err != nil {
		return nil, err
	}

	filter := model.ReviewFilter{
		CourseID: courseID,
		Page:     0,
		PageSize: 100,
	}

	reviews, err := GetReviewList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	llm, err := openai.New()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	inputJson, _ := json.Marshal(map[string]any{
		"courseName":    coursePO.Name,
		"teacherGroup":  offeredCoursePOs[courseID],
		"ratingAverage": infos[courseID].Average,
		"ratingCount":   infos[courseID].Count,
		"recentReviews": reviews,
	})

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, constant.GetCourseSummaryPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, string(inputJson)),
	}

	completion, err := llm.GenerateContent(
		context.Background(),
		content,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var response dto.GetCourseSummaryResponse
	err = json.Unmarshal([]byte(completion.Choices[0].Content), &response)

	return &response, err

}

// VectorizeCourse对课程进行向量化，
// 用于后续GetMatchCourses进行向量匹配。
// courseID为课程ID。
// TODO: 此处使用课程最近100条评论和课程名进行向量化，后续可以
// 优化和调整；
func VectorizeCourse(ctx context.Context, courseID int64) error {
	courseQuery := repository.NewCourseQuery(dal.GetDBClient())
	coursePOs, err := courseQuery.GetCourse(ctx, repository.WithID(courseID))
	if err != nil || len(coursePOs) == 0 {
		return err
	}
	coursePO := coursePOs[0]

	courseName := coursePO.Name

	filter := model.ReviewFilter{
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

// GetMatchCourses使用向量匹配，
// 根据自然语言描述找到最匹配的课程列表。
// description为用户提供的自然语言描述。
// TODO: 此处向量相似性计算（SimilaritySearch）中，
// 输出的课程列表数量为2，后续可以修改。

func GetMatchCourses(ctx context.Context, description string) ([]model.CourseSummary, error) {
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

	query := repository.NewCourseQuery(dal.GetDBClient())

	coursePOs, err := query.GetCourse(ctx, repository.WithIDs(courseIDs))
	if err != nil {
		return nil, err
	}

	courseCategories, err := query.GetCourseCategories(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	infos, err := ratingQuery.GetRatingInfoByIDs(ctx, model.RelatedTypeCourse, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]model.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithCategories(&course, courseCategories[int64(coursePO.ID)])
		converter.PackCourseWithRatingInfo(&course, infos[int64(coursePO.ID)])
		courses = append(courses, course)
	}
	return courses, nil

}