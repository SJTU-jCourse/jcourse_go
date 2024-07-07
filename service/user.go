package service

import (
	"context"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

func GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]domain.User, error) {
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

	result := make(map[int64]domain.User)
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
