package repository

import (
	"context"

	"gorm.io/gorm"
)

type DBOperation func(repo IRepository) error
type IRepository interface {
	NewUserPointQuery() IUserPointDetailQuery
	NewUserQuery() IUserQuery
	NewReviewQuery() IReviewQuery
	NewBaseCourseQuery() IBaseCourseQuery
	NewCourseQuery() ICourseQuery
	NewOfferedCourseQuery() IOfferedCourseQuery
	NewRatingQuery() IRatingQuery
	NewSettingQuery() ISettingQuery
	NewTeacherQuery() ITeacherQuery
	NewTrainingPlanQuery() ITrainingPlanQuery
	NewTrainingPlanCourseQuery() ITrainingPlanCourseQuery
	InTransaction(ctx context.Context, operation DBOperation) error
}
type Repository struct {
	db *gorm.DB
}

func (r *Repository) NewUserPointQuery() IUserPointDetailQuery {
	return &UserPointDetailQuery{
		db: r.db,
	}
}

func (r *Repository) NewReviewQuery() IReviewQuery {
	return &ReviewQuery{
		db: r.db,
	}
}

func (r *Repository) NewBaseCourseQuery() IBaseCourseQuery {
	return &BaseCourseQuery{
		db: r.db,
	}
}

func (r *Repository) NewCourseQuery() ICourseQuery {
	return &CourseQuery{
		db: r.db,
	}
}

func (r *Repository) NewOfferedCourseQuery() IOfferedCourseQuery {
	return &OfferedCourseQuery{
		db: r.db,
	}
}

func (r *Repository) NewRatingQuery() IRatingQuery {
	return &RatingQuery{
		db: r.db,
	}
}

func (r *Repository) NewSettingQuery() ISettingQuery {
	return &SettingQuery{
		db: r.db,
	}
}

func (r *Repository) NewTeacherQuery() ITeacherQuery {
	return &TeacherQuery{
		db: r.db,
	}
}

func (r *Repository) NewTrainingPlanQuery() ITrainingPlanQuery {
	return &TrainingPlanQuery{
		db: r.db,
	}
}

func (r *Repository) NewTrainingPlanCourseQuery() ITrainingPlanCourseQuery {
	return &TrainingPlanCourseQuery{
		db: r.db,
	}
}

var _ IRepository = &Repository{}

func (r *Repository) NewUserQuery() IUserQuery {
	return &UserQuery{
		db: r.db,
	}
}
func (r *Repository) NewUserPointDetailQuery() IUserPointDetailQuery {
	return &UserPointDetailQuery{
		db: r.db,
	}
}
func (r *Repository) InTransaction(ctx context.Context, operation DBOperation) error {
	db := r.db.WithContext(ctx)
	tx := db.Begin()
	if err := operation(NewRepository(tx)); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{db: db}
}
