package course

import (
	"errors"
	"time"

	"jcourse_go/internal/domain/shared"
)

type Review struct {
	ID        shared.IDType `json:"id"`
	CourseID  shared.IDType `json:"course_id"`
	UserID    shared.IDType `json:"user_id"`
	Comment   string        `json:"comment"`
	Rating    int64         `json:"rating"`
	Semester  string        `json:"semester"`
	IsPublic  bool          `json:"is_public"`
	Score     string        `json:"score"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
}

func (r *Review) Validate() error {
	if len(r.Comment) == 0 || len(r.Comment) > 10000 {
		return errors.New("comment length exceeds range [1, 10000]")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating is out of range [1,5]")
	}
	if len(r.Score) > 100 {
		return errors.New("score length exceeds range [1,100]")
	}
	return nil
}

func (r *Review) MakeRevision(updatedUserID shared.IDType, now time.Time) ReviewRevision {
	return ReviewRevision{
		ReviewID:  r.ID,
		Comment:   r.Comment,
		Rating:    r.Rating,
		Semester:  r.Semester,
		IsPublic:  r.IsPublic,
		Grade:     r.Score,
		CreatedAt: now,
		UpdatedBy: updatedUserID,
	}
}

func (r *Review) BeUpdated(cmd UpdateReviewCommand, now time.Time) error {
	r.Comment = cmd.Comment
	r.Rating = cmd.Rating
	r.UpdatedAt = now
	r.Semester = cmd.Semester
	r.Score = cmd.Score
	return r.Validate()
}

type ReviewRevision struct {
	ID       shared.IDType `json:"id"`
	ReviewID shared.IDType `json:"review_id"`
	Comment  string        `json:"comment"`
	Rating   int64         `json:"rating"`
	Semester string        `json:"semester"`
	IsPublic bool          `json:"is_public"`
	Grade    string        `json:"grade"`

	UpdatedBy shared.IDType `json:"updated_by"`
	CreatedAt time.Time     `json:"created_at"`
}

func NewReview(cmd WriteReviewCommand, userID shared.IDType, now time.Time) (Review, error) {
	review := Review{
		CourseID:  cmd.CourseID,
		UserID:    userID,
		Comment:   cmd.Comment,
		Rating:    cmd.Rating,
		Semester:  cmd.Semester,
		Score:     cmd.Score,
		IsPublic:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return review, review.Validate()
}
