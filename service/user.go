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

func GetUserSummaryByID(ctx context.Context, userID int64) (*dto.UserSummaryDTO, error) {
	filter := domain.ReviewFilter{
		UserID: userID,
	}

	total, _ := GetReviewCount(ctx, filter)
	// 过滤非匿名点评

	// 获取用户收到的赞数、被打赏积分数、关注的课程数

	return converter.ConvertUserDomainToUserSummaryDTO(userID, total, 0, 0, 0), nil
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

// 共用函数，用于获取用户基本信息和详细资料并组装成domain.User
func GetUserDomainByID(ctx context.Context, userID int64) (*domain.User, error) {
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

	user := converter.ConvertUserPOToDomain(*userPO)
	converter.PackUserWithProfile(&user, *userProfilePO)

	return &user, nil
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

func AdminGetUserList(ctx context.Context, filter domain.UserFilter) ([]dto.UserProfileDTO, error) {
	// 视前端而定获取用户的哪些信息
	// E.g. UserProfileDTO
	return nil, nil
}

func GetUserCount(ctx context.Context, filter domain.UserFilter) (int64, error) {
	userQuery := repository.NewUserQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	return userQuery.GetUserCount(ctx, opts...)
}

func UpdateUserProfileByID(ctx context.Context, userProfileDTO *dto.UserProfileDTO) error {
	userQuery := repository.NewUserQuery()
	oldUserPO, errQuery := userQuery.GetUserByID(ctx, userProfileDTO.UserID)
	if errQuery != nil {
		return errQuery
	}
	newUserPO := converter.ConvertUpdateUserProfileDTOToUserPO(userProfileDTO, oldUserPO)

	errUpdate := userQuery.UpdateUserByID(ctx, &newUserPO)
	if errUpdate != nil {
		return errUpdate
	}

	userProfileQuery := repository.NewUserProfileQuery()
	oldUserProfilePO, errQuery2 := userProfileQuery.GetUserProfileByID(ctx, userProfileDTO.UserID)
	if errQuery2 != nil {
		return errQuery2
	}
	newUserProfilePO := converter.ConvertUpdateUserProfileDTOToUsrProfilePO(userProfileDTO, oldUserProfilePO)
	errUpdate2 := userProfileQuery.UpdateUserProfileByID(ctx, &newUserProfilePO)
	if errUpdate2 != nil {
		return errUpdate2
	}
	return nil
}
