package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"jcourse_go/dal"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"jcourse_go/model/model"

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
	for _, po := range userPOs {
		result = append(result, converter.ConvertUserMinimalFromPO(po))
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
	if filter.StartTime != "" && filter.EndTime != "" {
		startTime, err := util.ParseTime(filter.StartTime)
		if err != nil {
			return opts
		}
		endTime, err := util.ParseTime(filter.EndTime)
		if err != nil {
			return opts
		}
		opts = append(opts, repository.WithTimeBetween(startTime, endTime))
	} else if filter.StartTime != "" {
		startTime, err := util.ParseTime(filter.StartTime)
		if err != nil {
			return opts
		}
		opts = append(opts, repository.WithTimeAfter(startTime))
	} else if filter.EndTime != "" {
		endTime, err := util.ParseTime(filter.EndTime)
		if err != nil {
			return opts
		}
		opts = append(opts, repository.WithTimeBefore(endTime))
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

func ChangeUserPoints(ctx context.Context, userID int64, eventType po.PointEventType, value int64, description string) error {
	userQuery := repository.NewUserQuery(dal.GetDBClient())
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
	user.Points += value
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	operation := func(db *gorm.DB) error {
		userQuery.UpdateUser(ctx, user)
		userPointDetailQuery.CreateUserPointDetail(ctx, userID, eventType, value, description)
		return nil
	}
	handler := repository.NewTransactionHandler(dal.GetDBClient())
	repository.InTransAction(ctx, handler, operation)
	return nil
}

const (
	handlerFeeRate = 0.01
)

func calcHandlingFee(value int64) int64 {
	return int64(float64(value) * (1 - handlerFeeRate))
}

func RedeemUserPoints(ctx context.Context, userID int64, value int64) error {
	// 给传承等开放的兑换积分接口
	userQuery := repository.NewUserQuery(dal.GetDBClient())
	userPOs, err := userQuery.GetUser(ctx, repository.WithID(userID))
	if err != nil {
		return err
	}
	if len(userPOs) == 0 {
		return errors.New("user not found")
	}
	user := userPOs[0]
	if user.Points < value {
		msg := fmt.Sprintf("user has not enough points, you have %d, require %d", user.Points, value)
		return errors.New(msg)
	}
	user.Points -= value
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	operation := func(db *gorm.DB) error {
		userQuery.UpdateUser(ctx, user)
		userPointDetailQuery.CreateUserPointDetail(ctx, userID, po.PointEventRedeem, -value, fmt.Sprintf("用户%d兑换积分%d", userID, value))
		return nil
	}
	handler := repository.NewTransactionHandler(dal.GetDBClient())
	err = repository.InTransAction(ctx, handler, operation)
	return err
}

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
	receivedValue := value - calcHandlingFee(value)
	senderPO.Points -= value
	receiverPO.Points += receivedValue
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	handler := repository.NewTransactionHandler(dal.GetDBClient())
	description := fmt.Sprintf("用户%d转账给用户%d %d分", senderID, receiverID, value)
	operation := func(db *gorm.DB) error {
		userQuery.UpdateUser(ctx, *senderPO)
		userQuery.UpdateUser(ctx, *receiverPO)
		userPointDetailQuery.CreateUserPointDetail(ctx, senderID, po.PointEventTransfer, -value, description)
		userPointDetailQuery.CreateUserPointDetail(ctx, receiverID, po.PointEventTransfer, receivedValue, description)
		return nil
	}
	err = repository.InTransAction(ctx, handler, operation)
	return err
}
