package service

import (
	"context"
	"jcourse_go/model/dto"
	"jcourse_go/util"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

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
	userSummary := converter.ConvertToUserSummaryDTO(userPO, userProfilePO)
	return userSummary, nil
}

func GetUserDetailsByID(ctx context.Context, userID int64) (*dto.UserDetailsDTO, error) {
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
	userDetails := converter.ConvertToUserDetailsDTO(userPO, userProfilePO)
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

func GetUserList(ctx context.Context, filter domain.UserFilter) ([]dto.UserSummaryDTO, error) {
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
	result := make([]dto.UserSummaryDTO, 0)

	userProfileMap := make(map[int]string)
	for _, profile := range userProfilePOs {
		userProfileMap[int(profile.UserID)] = profile.Avatar
	}

	for _, userPO := range userPOs {
		if avatar, exists := userProfileMap[int(userPO.ID)]; exists {
			userSummaryDTO := dto.UserSummaryDTO{
				ID:       int64(userPO.ID),
				Username: userPO.Username,
				Role:     userPO.UserRole,
				Avatar:   avatar,
			}
			result = append(result, userSummaryDTO)
		} else {
			userSummaryDTO := dto.UserSummaryDTO{
				ID:       int64(userPO.ID),
				Username: userPO.Username,
				Role:     userPO.UserRole,
				Avatar:   "",
			}
			result = append(result, userSummaryDTO)
		}
	}
	return result, nil
}

func GetUserCount(ctx context.Context, filter domain.UserFilter) (int64, error) {
	userQuery := repository.NewUserQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserDBOptionFromFilter(userQuery, filter)
	return userQuery.GetUserCount(ctx, opts...)
}
