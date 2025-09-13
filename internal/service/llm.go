package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/review"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/infrastructure/rpc"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/model/converter"
	"jcourse_go/internal/model/types"
)

var llm *openai.LLM

func InitLLM() error {
	var err error
	llm, err = openai.New()
	if err != nil {
		return err
	}
	return nil
}

// OptCourseReview使用LLM提示词对课程评价内容进行优化。
// courseName为课程名称，
// reviewContent为评价的内容，
// 函数返回值包含两个字段：
// Suggestion为修改建议，
// Result为根据修改建议给出的一种修改结果。
func OptCourseReview(ctx context.Context, courseName string, reviewContent string) (olddto.OptCourseReviewResponse, error) {
	inputJson, _ := json.Marshal(map[string]string{
		"course": courseName,
		"review": reviewContent,
	})

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, constant.OptCourseReviewPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, string(inputJson)),
	}

	completion, err := llm.GenerateContent(
		ctx,
		content,
	)

	if err != nil {
		return olddto.OptCourseReviewResponse{}, err
	}
	var response olddto.OptCourseReviewResponse
	responseStr := trimLLMJSON(completion.Choices[0].Content)
	err = json.Unmarshal([]byte(responseStr), &response)

	return response, err
}

func trimLLMJSON(raw string) string {
	s := strings.TrimPrefix(raw, "```json")
	s = strings.TrimSuffix(s, "```")
	return s
}

func buildCourseSummaryPrompt(course *course.CourseDetail, reviews []course.Review) string {
	type tinyReview struct {
		Rating  int64  `json:"rating,omitempty"`
		Comment string `json:"comment,omitempty"`
	}
	tinyReviews := make([]tinyReview, 0)

	for _, r := range reviews {
		tinyReviews = append(tinyReviews, tinyReview{
			Rating:  r.Rating,
			Comment: r.Comment,
		})
	}

	coursePrompt := struct {
		CourseName    string       `json:"course_name,omitempty"`
		TeacherName   string       `json:"teacher_name,omitempty"`
		RatingAverage float64      `json:"rating_average,omitempty"`
		RatingCount   int64        `json:"rating_count,omitempty"`
		RecentReviews []tinyReview `json:"recent_reviews,omitempty"`
	}{
		CourseName:    course.Name,
		TeacherName:   course.MainTeacher.Name,
		RatingAverage: course.RatingInfo.Average,
		RatingCount:   course.RatingInfo.Count,
		RecentReviews: tinyReviews,
	}

	p, _ := json.Marshal(coursePrompt)
	return string(p)
}

// GetCourseSummary使用LLM提示词基于课程评价生成课程总结。
// courseID为课程的ID，
// 返回值包含课程的总结。
// TODO: 此处基于课程最近的100条评价内容生成课程总结，后续可
// 以进一步调整和优化。
func GetCourseSummary(ctx context.Context, courseID int64) (*olddto.GetCourseSummaryResponse, error) {
	courseDetail, err := GetCourseDetail(ctx, courseID, 0)
	if err != nil {
		return nil, err
	}

	filter := reaction.ReviewFilterForQuery{
		CourseID: courseID,
		PaginationFilterForQuery: course.PaginationFilterForQuery{
			Page:     0,
			PageSize: 100,
		},
	}

	reviews, err := GetReviewList(ctx, nil, filter)
	if err != nil {
		return nil, err
	}

	inputJson := buildCourseSummaryPrompt(courseDetail, reviews)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, constant.GetCourseSummaryPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, inputJson),
	}

	completion, err := llm.GenerateContent(
		ctx,
		content,
	)

	if err != nil {
		return nil, err
	}

	var response olddto.GetCourseSummaryResponse
	responseStr := trimLLMJSON(completion.Choices[0].Content)
	err = json.Unmarshal([]byte(responseStr), &response)

	return &response, err

}

// VectorizeCourse对课程进行向量化，
// 用于后续GetMatchCourses进行向量匹配。
// courseID为课程ID。
// TODO: 此处使用课程最近100条评论和课程名进行向量化，后续可以
// 优化和调整；
func VectorizeCourse(ctx context.Context, courseID int64) error {
	c := repository.Q.CoursePO
	coursePO, err := c.WithContext(ctx).Where(c.ID.Eq(courseID)).Take()
	if err != nil {
		return err
	}

	courseName := coursePO.Name

	filter := reaction.ReviewFilterForQuery{
		CourseID: courseID,
		PaginationFilterForQuery: course.PaginationFilterForQuery{
			Page:     0,
			PageSize: 100,
		},
	}

	reviews, err := GetReviewList(ctx, nil, filter)
	if err != nil {

		return err
	}

	var comments []string

	for _, review := range reviews {
		comments = append(comments, review.Comment)
	}

	vectorStore, err := rpc.OpenVectorStoreConn(ctx)

	if err != nil {

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
		ctx,
		[]schema.Document{doc},
	)

	if err != nil {

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

func GetMatchCourses(ctx context.Context, description string) ([]course.CourseSummary, error) {
	vectorStore, err := rpc.OpenVectorStoreConn(ctx)

	if err != nil {

		return nil, err
	}

	docs, err := vectorStore.SimilaritySearch(ctx, description, 2)
	if err != nil {

		return nil, err
	}

	err = vectorStore.Close()
	if err != nil {

		return nil, err
	}

	var courseIDs []int64
	for _, doc := range docs {
		courseID := doc.Metadata["courseID"].(float64)
		courseIDs = append(courseIDs, int64(courseID))
	}

	c := repository.Q.CoursePO
	coursePOs, err := c.WithContext(ctx).Preload(c.Categories).Where(c.ID.In(courseIDs...)).Find()
	if err != nil {
		return nil, err
	}

	infos, err := GetMultipleRating(ctx, types.RelatedTypeCourse, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]course.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithRatingInfo(&course, infos[coursePO.ID])
		courses = append(courses, course)
	}
	return courses, nil

}
