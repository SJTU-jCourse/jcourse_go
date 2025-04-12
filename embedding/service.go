package embedding

import (
	"context"
	"fmt"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/types"
	"jcourse_go/repository"
	"jcourse_go/service"
	"log"
	"strings"

	"github.com/tmc/langchaingo/schema"
)

func VectorizeCourse(ctx context.Context, courseID int64) error {
	c := repository.Q.CoursePO
	coursePO, err := c.WithContext(ctx).Where(c.ID.Eq(courseID)).Take()
	if err != nil {
		return err
	}

	courseName := coursePO.Name

	filter := model.ReviewFilterForQuery{
		CourseID: courseID,
		PaginationFilterForQuery: model.PaginationFilterForQuery{
			Page:     0,
			PageSize: 100,
		},
	}

	reviews, err := service.GetReviewList(ctx, nil, filter)
	if err != nil {

		return err
	}

	var comments []string

	for _, review := range reviews {
		comments = append(comments, review.Comment)
	}

	vectorStore := GetStore()

	targetStr := courseName
	if len(comments) > 0 {
		targetStr += "\n" + strings.Join(comments, "\n")
	}

	doc := schema.Document{
		PageContent: targetStr,
		Metadata: map[string]any{
			"courseID": courseID,
		},
	}
	log.Printf("Vectorizing course %d using OpenAI", courseID)

	_, err = vectorStore.AddDocuments(ctx, []schema.Document{doc})
	if err != nil {
		log.Printf("Error adding document for course %d to vector store: %v", courseID, err)
		return fmt.Errorf("failed to vectorize course %d: %w", courseID, err)
	}

	log.Printf("Successfully vectorized course %d", courseID)
	return nil
}

func GetMatchCourses(ctx context.Context, description string) ([]model.CourseSummary, error) {
	vectorStore := GetStore()

	log.Printf("Performing similarity search using OpenAI for: %q", description)
	docs, err := vectorStore.SimilaritySearch(ctx, description, 5) // Top 5 results
	if err != nil {
		log.Printf("Error during similarity search: %v", err)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(docs) == 0 {
		log.Println("No matching courses found.")
		return []model.CourseSummary{}, nil
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

	infos, err := service.GetMultipleRating(ctx, types.RelatedTypeCourse, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]model.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithRatingInfo(&course, infos[coursePO.ID])
		courses = append(courses, course)
	}
	return courses, nil
}
