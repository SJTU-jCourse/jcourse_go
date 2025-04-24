package service

import (
	"context"
	"errors"

	"gorm.io/gorm/clause"

	"jcourse_go/internal/infra/query"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

func CreateRating(ctx context.Context, userID int64, dto dto.RatingDTO) error {

	if !types.IsARatingRelatedType(dto.RelatedType) {
		return errors.New("invalid related type")
	}

	ratingPO := converter.ConvertRatingDTOToPO(userID, dto)
	r := query.Q.RatingPO
	err := r.WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
			Columns: []clause.Column{
				{Name: string(r.UserID.ColumnName())},
				{Name: string(r.RelatedID.ColumnName())},
				{Name: string(r.RelatedType.ColumnName())},
			}}).
		Create(&ratingPO)
	if err != nil {
		return err
	}
	err = SyncRating(ctx, types.RatingRelatedType(ratingPO.RelatedType), ratingPO.RelatedID)
	if err != nil {
		return err
	}
	return nil
}

func GetRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64) (model.RatingInfo, error) {
	res := model.RatingInfo{}
	if !types.IsARatingRelatedType(string(relatedType)) {
		return res, errors.New("invalid related type")
	}

	dist := make([]model.RatingInfoDistItem, 0)

	r := query.Q.RatingPO
	err := r.WithContext(ctx).Select(r.Rating, r.ID.Count().As("count")).
		Where(r.RelatedID.Eq(relatedID), r.RelatedType.Eq(string(relatedType))).
		Group(r.Rating).Scan(&dist)
	if err != nil {
		return res, err
	}
	res.RatingDist = dist
	res.Calc()
	return res, nil
}

func GetMultipleRating(ctx context.Context, relatedType types.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error) {
	res := make(map[int64]model.RatingInfo)
	dist := make(map[int64][]model.RatingInfoDistItem)
	if !types.IsARatingRelatedType(string(relatedType)) {
		return res, errors.New("invalid related type")
	}

	r := query.Q.RatingPO
	rows := make([]struct {
		RelatedID int64 `json:"related_id"`
		Rating    int64 `json:"rating"`
		Count     int64 `json:"count"`
	}, 0)
	err := r.WithContext(ctx).Select(r.RelatedID, r.Rating, r.ID.Count().As("count")).
		Where(r.RelatedID.In(relatedIDs...), r.RelatedType.Eq(string(relatedType))).
		Group(r.RelatedID, r.Rating).Scan(&rows)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		dist[row.RelatedID] = append(dist[row.RelatedID], model.RatingInfoDistItem{
			Rating: row.Rating,
			Count:  row.Count,
		})
	}

	for id, distItems := range dist {
		ratingInfo := model.RatingInfo{RatingDist: distItems}
		ratingInfo.Calc()
		res[id] = ratingInfo
	}

	return res, nil
}

func GetUserRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64, userID int64) (int64, error) {
	if !types.IsARatingRelatedType(string(relatedType)) {
		return 0, errors.New("invalid related type")
	}

	r := query.Q.RatingPO
	rating, err := r.WithContext(ctx).Select(r.Rating).Where(r.RelatedID.Eq(relatedID), r.RelatedType.Eq(string(relatedType)), r.UserID.Eq(userID)).Take()
	if err != nil {
		return 0, err
	}
	return rating.Rating, nil
}

func SyncRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64) error {
	if !types.IsARatingRelatedType(string(relatedType)) {
		return errors.New("invalid related type")
	}

	ratingInfo, err := GetRating(ctx, relatedType, relatedID)
	if err != nil {
		return err
	}

	if relatedType == types.RelatedTypeCourse {
		return SyncCourseRating(ctx, relatedID, ratingInfo)
	} else if relatedType == types.RelatedTypeTeacher {
		return SyncTeacherRating(ctx, relatedID, ratingInfo)
	} else if relatedType == types.RelatedTypeTrainingPlan {
		return SyncTrainingPlanRating(ctx, relatedID, ratingInfo)
	}

	return nil
}

func SyncCourseRating(ctx context.Context, courseID int64, ratingInfo model.RatingInfo) error {
	c := query.Q.CoursePO
	_, err := c.WithContext(ctx).Select(c.RatingCount, c.RatingAvg).Where(c.ID.Eq(courseID)).
		Updates(po.CoursePO{RatingCount: ratingInfo.Count, RatingAvg: ratingInfo.Average})
	return err
}

func SyncTeacherRating(ctx context.Context, teacherID int64, ratingInfo model.RatingInfo) error {
	t := query.Q.TeacherPO
	_, err := t.WithContext(ctx).Select(t.RatingCount, t.RatingAvg).Where(t.ID.Eq(teacherID)).
		Updates(po.CoursePO{RatingCount: ratingInfo.Count, RatingAvg: ratingInfo.Average})
	return err
}

func SyncTrainingPlanRating(ctx context.Context, trainingPlanID int64, ratingInfo model.RatingInfo) error {
	tp := query.Q.TrainingPlanPO
	_, err := tp.WithContext(ctx).Select(tp.RatingCount, tp.RatingAvg).Where(tp.ID.Eq(trainingPlanID)).
		Updates(po.CoursePO{RatingCount: ratingInfo.Count, RatingAvg: ratingInfo.Average})
	return err
}
