package service

import (
	"context"

	"jcourse_go/internal/model/converter"
	"jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	repository2 "jcourse_go/internal/repository"
	"jcourse_go/pkg/util"
)

func GetUserActivityByID(ctx context.Context, userID int64) (*model.UserActivity, error) {
	// filter := model.ReviewFilterForQuery{
	// 	UserID: userID,
	// }

	// total, _ := GetReviewCount(ctx, filter)
	// 过滤非匿名点评

	// 获取用户收到的赞数、被打赏积分数、关注的课程数

	return &model.UserActivity{}, nil
}

func GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]model.UserMinimal, error) {
	result := make(map[int64]model.UserMinimal)
	if len(userIDs) == 0 {
		return result, nil
	}
	u := repository2.Q.UserPO
	userPOs, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return result, err
	}
	for _, userPO := range userPOs {
		user := converter.ConvertUserMinimalFromPO(userPO)
		result[user.ID] = user
	}
	return result, nil
}

// 共用函数，用于获取用户基本信息和详细资料并组装成domain.User
func GetUserDetailByID(ctx context.Context, userID int64) (*model.UserDetail, error) {
	u := repository2.Q.UserPO
	userPO, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).Take()
	if err != nil {
		return nil, err
	}
	user := converter.ConvertUserDetailFromPO(userPO)
	return &user, nil
}

func buildUserDBOptionFromFilter(ctx context.Context, q *repository2.Query, filter model.UserFilterForQuery) repository2.IUserPODo {
	builder := q.UserPO.WithContext(ctx)
	u := q.UserPO
	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.Order != "" {
		field, ok := u.GetFieldByName(filter.Order)
		if ok {
			if filter.Ascending {
				builder = builder.Order(field)
			} else {
				builder = builder.Order(field.Desc())
			}
		}
	}
	return builder
}

func GetUserList(ctx context.Context, filter model.UserFilterForQuery) ([]model.UserMinimal, error) {
	q := buildUserDBOptionFromFilter(ctx, repository2.Q, filter)
	userPOs, err := q.Find()
	if err != nil {
		return nil, err
	}

	result := make([]model.UserMinimal, 0)
	for _, userPO := range userPOs {
		result = append(result, converter.ConvertUserMinimalFromPO(userPO))
	}
	return result, nil
}

func GetUserCount(ctx context.Context, filter model.UserFilterForQuery) (int64, error) {
	filter.Page, filter.PageSize = 0, 0
	q := buildUserDBOptionFromFilter(ctx, repository2.Q, filter)
	return q.Count()
}

func UpdateUserProfileByID(ctx context.Context, userProfileDTO dto.UserProfileDTO, userID int64) error {
	u := repository2.Q.UserPO
	newUserPO := converter.ConvertUserProfileToPO(userProfileDTO)
	newUserPO.ID = userID
	err := u.WithContext(ctx).Save(&newUserPO)
	if err != nil {
		return err
	}
	return nil
}
