package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type IRatingQuery interface {
	GetRatingInfo(ctx context.Context, relatedType model.RatingRelatedType, relatedID int64) (model.RatingInfo, error)
	GetUserRating(ctx context.Context, relatedType model.RatingRelatedType, relatedID int64, userID int64) (int64, error)
	GetRatingInfoByIDs(ctx context.Context, relatedType model.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error)
	CreateRating(ctx context.Context, ratingPO po.RatingPO) error
	UpdateRating(ctx context.Context, ratingPO po.RatingPO) error
	DeleteRating(ctx context.Context, ratingPO po.RatingPO) error
	GetRating(ctx context.Context, userID int64, relatedID int64, relatedType model.RatingRelatedType) (po.RatingPO, error)
}

type RatingQuery struct {
	db *gorm.DB
}

func (r *RatingQuery) GetRating(ctx context.Context, userID int64, relatedID int64, relatedType model.RatingRelatedType) (po.RatingPO, error) {
	db := r.optionDB(ctx)
	ratingPO := po.RatingPO{}
	result := db.Where("user_id = ? and related_id = ? and related_type = ?", userID, relatedID, relatedType).Take(&ratingPO)
	return ratingPO, result.Error
}

func (r *RatingQuery) GetUserRating(ctx context.Context, relatedType model.RatingRelatedType, relatedID int64, userID int64) (int64, error) {
	db := r.optionDB(ctx)
	ratingPO := po.RatingPO{}
	result := db.Where("user_id = ? and related_id = ? and related_type = ?", userID, relatedID, relatedType).Take(&ratingPO)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	}
	return ratingPO.Rating, nil
}

func (r *RatingQuery) UpdateRating(ctx context.Context, ratingPO po.RatingPO) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&ratingPO).Error; err != nil {
			return err
		}
		if err := r.SyncRating(ctx, tx, ratingPO.RelatedID, ratingPO.RelatedType); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RatingQuery) DeleteRating(ctx context.Context, ratingPO po.RatingPO) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&po.RatingPO{}).Where("user_id = ? and related_type = ? and related_id = ?", ratingPO.UserID, ratingPO.RelatedType, ratingPO.RelatedID).Delete(&ratingPO).Error; err != nil {
			return err
		}
		if err := r.SyncRating(ctx, tx, ratingPO.RelatedID, ratingPO.RelatedType); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RatingQuery) CreateRating(ctx context.Context, ratingPO po.RatingPO) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&po.RatingPO{}).Create(&ratingPO).Error; err != nil {
			return err
		}
		if err := r.SyncRating(ctx, tx, ratingPO.RelatedID, ratingPO.RelatedType); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RatingQuery) SyncRating(ctx context.Context, db *gorm.DB, relatedID int64, relatedType string) error {
	// 1. Aggregate review data
	var info po.RatingInfo
	if err := db.WithContext(ctx).
		Model(&po.RatingPO{}).
		Select("COUNT(*) AS count, AVG(rating) AS average").
		Where("related_type = ? and related_id = ?", relatedType, relatedID).
		Scan(&info).
		Error; err != nil {
		return err
	}

	var targetModelMap = map[model.RatingRelatedType]any{
		model.RelatedTypeCourse:       &po.CoursePO{},
		model.RelatedTypeTeacher:      &po.TeacherPO{},
		model.RelatedTypeTrainingPlan: &po.TrainingPlanPO{},
	}

	targetModel := targetModelMap[relatedType]

	// 2. Update the matching course row
	if err := db.WithContext(ctx).
		Model(targetModel).
		Where("id = ?", relatedID).
		Updates(map[string]interface{}{
			"rating_count": info.Count,
			"rating_avg":   info.Average,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *RatingQuery) GetRatingInfoByIDs(ctx context.Context, relatedType model.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error) {
	res := make(map[int64]model.RatingInfo)
	distByIDs := make([]model.RatingInfoDistItemByID, 0)
	db := r.optionDB(ctx)
	err := db.Select("rating, count(id) as count, related_id").
		Where("related_type = ? and related_id in ?", relatedType, relatedIDs).
		Group("rating").Group("related_id").
		Find(&distByIDs).Error
	if err != nil {
		return nil, err
	}
	for _, dist := range distByIDs {
		info, ok := res[dist.RelatedID]
		if !ok {
			info = model.RatingInfo{RatingDist: make([]model.RatingInfoDistItem, 0)}
		}
		info.RatingDist = append(info.RatingDist, model.RatingInfoDistItem{Rating: dist.Rating, Count: dist.Count})
		info.Calc()
		res[dist.RelatedID] = info
	}
	return res, nil
}

func (r *RatingQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&po.RatingPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (r *RatingQuery) GetRatingInfo(ctx context.Context, relatedType model.RatingRelatedType, relatedID int64) (model.RatingInfo, error) {
	res := model.RatingInfo{}
	dists := make([]model.RatingInfoDistItem, 0)
	db := r.optionDB(ctx)
	err := db.Select("rating, count(*) as count").
		Where("related_type = ? and related_id = ?", relatedType, relatedID).
		Group("rating").
		Find(&dists).Error
	if err != nil {
		return res, err
	}
	res.RatingDist = dists
	res.Calc()
	return res, nil
}

func NewRatingQuery(db *gorm.DB) IRatingQuery {
	return &RatingQuery{db: db}
}
