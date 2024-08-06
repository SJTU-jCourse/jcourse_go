package service

import (
	"context"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

// 该函数如果放在converter包中会报错import cycle is not allowed
// UserSummaryDTO的组装需要借助service/review.go中的GetReviewCount函数
func ConvertToUserSummaryDTO(ctx context.Context, userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserSummaryDTO {
	if userPO == nil {
		return nil
	}

	filter := domain.ReviewFilter{
		UserID: int64(userPO.ID),
	}

	total, _ := GetReviewCount(ctx, filter)

	return &dto.UserSummaryDTO{
		ID:                   int64(userPO.ID),
		ReviewCount:          total,
		TipReceive:           0,
		FollowingCourseCount: 0,
	}
}

func GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]domain.User, error) {
	result := make(map[int64]domain.User)
	if len(userIDs) == 0 {
		return result, nil
	}

	userQuery := repository.NewUserQuery()
	userMap, err := userQuery.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userProfileQuery := repository.NewUserProfileQuery()
	userProfileMap, err := userProfileQuery.GetUserProfileByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for _, userPO := range userMap {
		user := converter.ConvertUserPOToDomain(userPO)
		profilePO, ok := userProfileMap[user.ID]
		if ok {
			converter.PackUserWithProfile(&user, profilePO)
		}
		result[user.ID] = user
	}
	return result, nil
}

func GetUserSummaryByID(ctx context.Context, userID int64) (*dto.UserSummaryDTO, error) {
	userQuery := repository.NewUserQuery()
	userPO, err := userQuery.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userProfileQuery := repository.NewUserProfileQuery()
	userProfilePO, err := userProfileQuery.GetUserProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userSummary := ConvertToUserSummaryDTO(ctx, userPO, userProfilePO)
	return userSummary, nil
}

func GetUserDetailByID(ctx context.Context, userID int64) (*dto.UserDetailDTO, error) {
	userQuery := repository.NewUserQuery()
	userPO, err := userQuery.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userProfileQuery := repository.NewUserProfileQuery()
	userProfilePO, err := userProfileQuery.GetUserProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userDetails := converter.ConvertToUserDetailDTO(userPO, userProfilePO)
	return userDetails, nil
}

func GetUserProfileByID(ctx context.Context, userID int64) (*dto.UserProfileDTO, error) {
	userQuery := repository.NewUserQuery()
	userPO, err := userQuery.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userProfileQuery := repository.NewUserProfileQuery()
	userProfilePO, err := userProfileQuery.GetUserProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userProfile := converter.ConvertToUserProfileDTO(userPO, userProfilePO)
	return userProfile, nil
}

func buildUserDBOptionFromFilter(query repository.IUserQuery, filter domain.UserFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.PageSize > 0 {
		opts = append(opts, query.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, query.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	return opts
}

func GetUserList(ctx context.Context, filter domain.UserFilter) ([]dto.UserDetailDTO, error) {
	userQuery := repository.NewUserQuery()
	userProfileQuery := repository.NewUserProfileQuery()
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	userPOs, err := userQuery.GetUserList(ctx, opts...)
	if err != nil {
		return nil, err
	}
	userProfilePOs, err := userProfileQuery.GetUserProfileList(ctx, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]dto.UserDetailDTO, 0)

	userProfileMap := make(map[int]*po.UserProfilePO)
	for _, userProfilePO := range userProfilePOs {
		userProfileMap[int(userProfilePO.UserID)] = &userProfilePO
	}

	for _, userPO := range userPOs {
		userDetailDTO := dto.UserDetailDTO{
			ID:       int64(userPO.ID),
			Username: userPO.Username,
			Avatar:   "",
			Bio:      "",
		}
		if userProfilePO, exists := userProfileMap[int(userPO.ID)]; exists {
			userDetailDTO.Avatar = userProfilePO.Avatar
			userDetailDTO.Bio = userProfilePO.Bio
		}
		result = append(result, userDetailDTO)
	}
	return result, nil
}

func GetUserCount(ctx context.Context, filter domain.UserFilter) (int64, error) {
	userQuery := repository.NewUserQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	return userQuery.GetUserCount(ctx, opts...)
}
