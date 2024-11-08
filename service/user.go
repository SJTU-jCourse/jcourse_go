package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"jcourse_go/model/converter"
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
	if filter.PageSize > 0 {
		opts = append(opts, repository.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, repository.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	if !filter.StartTime.IsZero() && !filter.EndTime.IsZero() {
		opts = append(opts, repository.WithTimeBetween(filter.StartTime, filter.EndTime))
	} else if !filter.StartTime.IsZero() {
		opts = append(opts, repository.WithTimeAfter(filter.StartTime))
	} else if !filter.EndTime.IsZero() {
		opts = append(opts, repository.WithTimeBefore(filter.EndTime))
	}
	return opts
}

func GetUserPointDetailList(ctx context.Context, fileter model.UserPointDetailFilter) ([]model.UserPointDetailItem, error) {
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	opts := buildUserPointDetailDBOptionFromFilter(userPointDetailQuery, fileter)
	userPointDetailPOs, err := userPointDetailQuery.GetUserPointDetail(ctx, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]model.UserPointDetailItem, 0)
	for _, detailPO := range userPointDetailPOs {
		result = append(result, converter.ConvertUserPointDetailItemFromPO(detailPO))
	}
	return result, nil
}
func GetUserPointDetailCount(ctx context.Context, filter model.UserPointDetailFilter) (int64, error) {
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserPointDetailDBOptionFromFilter(userPointDetailQuery, filter)
	return userPointDetailQuery.GetUserPointDetailCount(ctx, opts...)
}

// HINT: 以下的几个UserPoint相关函数都是并发安全的, 但不保证成功，事务失败时需要上层自行处理
func ChangeUserPoints(ctx context.Context, userID int64, eventType model.PointEventType, value int64, description string) error {
	repo := repository.NewRepository(dal.GetDBClient())
	userQuery := repo.NewUserQuery()
	userPOs, err := userQuery.GetUser(ctx, repository.WithID(userID))
	if err != nil {
		return err
	}
	if len(userPOs) == 0 {
		return errors.Errorf("user %d not found", userID)
	}
	user := userPOs[0]
	if user.Points+value < 0 {
		return errors.Errorf("user %d has not enough points", userID)
	}
	originalPoints := user.Points
	user.Points += value
	userPointDetailQuery := repo.NewUserPointQuery()
	operation := func(repo repository.IRepository) error {
		userQuery.UpdateUser(ctx, user, repository.WithOptimisticLock("points", originalPoints))
		userPointDetailQuery.CreateUserPointDetail(ctx, userID, eventType, value, description)
		return nil
	}
	return repo.InTransaction(ctx, operation)
}

func calcHandlingFee(ctx context.Context, value int64) int64 {
	siteQuery := repository.NewSettingQuery(dal.GetDBClient())
	siteSetting, err := siteQuery.GetSetting(ctx, constant.HandleFeeRateKey)
	if err != nil || siteSetting == nil {
		return int64(float64(value) * (1 - constant.DefaultHandleFeeRate))
	}
	handlerFeeRate := siteSetting.GetValue().(float64)
	return int64(float64(value) * (1 - handlerFeeRate))
}

const (
	TransferDescriptionFormat = "用户%d转账给用户%d %d分"
)

func TransferUserPoints(ctx context.Context, senderID int64, receiverID int64, value int64) error {
	userQuery := repository.NewUserQuery(dal.GetDBClient())

	// 合并到一次查询
	ids := []int64{senderID, receiverID}
	userPOs, err := userQuery.GetUser(ctx, repository.WithIDs(ids))
	if err != nil {
		return err
	}
	var senderPO *po.UserPO = nil
	var receiverPO *po.UserPO = nil
	for _, user := range userPOs {
		if user.ID == uint(senderID) {
			senderPO = &user
		}
		if user.ID == uint(receiverID) {
			receiverPO = &user
		}
	}
	if senderPO == nil {
		return errors.New("sender not found")
	}
	if receiverPO == nil {
		return errors.New("receiver not found")
	}
	if senderPO.Points < value {
		return errors.New("sender has not enough points")
	}
	receivedValue := value - calcHandlingFee(ctx, value)
	senderOriginalPoints := senderPO.Points
	receiverOriginalPoints := receiverPO.Points
	senderPO.Points -= value
	receiverPO.Points += receivedValue
	repo := repository.NewRepository(dal.GetDBClient())
	description := fmt.Sprintf(TransferDescriptionFormat, senderID, receiverID, value)
	var operations repository.DBOperation
	operations = func(repo repository.IRepository) error {
		userQuery := repo.NewUserQuery()
		userPointDetailQuery := repo.NewUserPointQuery()
		userQuery.UpdateUser(ctx, *senderPO, repository.WithOptimisticLock("points", senderOriginalPoints))
		userQuery.UpdateUser(ctx, *receiverPO, repository.WithOptimisticLock("points", receiverOriginalPoints))
		userPointDetailQuery.CreateUserPointDetail(ctx, senderID, model.PointEventTransfer, -value, description)
		userPointDetailQuery.CreateUserPointDetail(ctx, receiverID, model.PointEventTransfer, receivedValue, description)
		return nil
	}
	err = repo.InTransaction(ctx, operations)
	return err
}
