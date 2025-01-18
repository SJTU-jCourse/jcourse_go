package service

import (
	"context"
	"time"

	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/repository"
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

	userQuery := repository.NewUserQuery(dal.GetDBClient())
	userMap, err := userQuery.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for _, userPO := range userMap {
		user := converter.ConvertUserMinimalFromPO(userPO)
		result[user.ID] = user
	}
	return result, nil
}

// 共用函数，用于获取用户基本信息和详细资料并组装成domain.User
func GetUserDetailByID(ctx context.Context, userID int64) (*model.UserDetail, error) {
	userQuery := repository.NewUserQuery(dal.GetDBClient())
	userPO, err := userQuery.GetUser(ctx, repository.WithID(userID))
	if err != nil || len(userPO) == 0 {
		return nil, err
	}
	user := converter.ConvertUserDetailFromPO(userPO[0])
	return &user, nil
}

func buildUserDBOptionFromFilter(query repository.IUserQuery, filter model.UserFilterForQuery) []repository.DBOption {
	opts := buildPaginationDBOptions(filter.PaginationFilterForQuery)
	return opts
}

func GetUserList(ctx context.Context, filter model.UserFilterForQuery) ([]model.UserMinimal, error) {
	userQuery := repository.NewUserQuery(dal.GetDBClient())
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	userPOs, err := userQuery.GetUser(ctx, opts...)
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
	userQuery := repository.NewUserQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	return userQuery.GetUserCount(ctx, opts...)
}

func UpdateUserProfileByID(ctx context.Context, userProfileDTO dto.UserProfileDTO, userID int64) error {
	userQuery := repository.NewUserQuery(dal.GetDBClient())
	newUserPO := converter.ConvertUserProfileToPO(userProfileDTO)
	newUserPO.ID = uint(userID)
	errUpdate := userQuery.UpdateUser(ctx, newUserPO)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

func buildUserPointDetailDBOptionFromFilter(query repository.IUserPointDetailQuery, filter model.UserPointDetailFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.UserPointDetailID > 0 {
		opts = append(opts, repository.WithID(filter.UserPointDetailID))
	}
	if filter.UserID > 0 {
		opts = append(opts, repository.WithUserID(filter.UserID))
	}
	if filter.Page > 0 && filter.PageSize > 0 {
		opts = append(opts, repository.WithPaginate(filter.Page, filter.PageSize))
	}
	if filter.StartTime > 0 && filter.EndTime > 0 {
		opts = append(opts, repository.WithCreatedAtBetween(time.Unix(filter.StartTime, 0), time.Unix(filter.EndTime, 0)))
	} else if filter.StartTime > 0 {
		opts = append(opts, repository.WithCreatedAtAfter(time.Unix(filter.StartTime, 0)))
	} else if filter.EndTime > 0 {
		opts = append(opts, repository.WithCreatedAtBefore(time.Unix(filter.EndTime, 0)))
	}
	return opts
}
