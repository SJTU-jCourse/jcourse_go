package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type IRatingQuery interface {
	GetRatingInfo(ctx context.Context, relatedType model.RatingRelatedType, relatedID int64) (model.RatingInfo, error)
	GetRatingInfoByIDs(ctx context.Context, relatedType model.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error)
}

type RatingQuery struct {
	db *gorm.DB
}

func (r *RatingQuery) GetRatingInfoByIDs(ctx context.Context, relatedType model.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error) {
	res := make(map[int64]model.RatingInfo)
	distByIDs := make([]model.RatingInfoDistItemByID, 0)
	db := r.optionDB(ctx)
	err := db.Select("rating, count(id)").
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
	err := db.Select("rating, count(id)").
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
