package service

import (
	"context"
	"fmt"
	"time"

	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
	"jcourse_go/repository"

	"github.com/pkg/errors"
)

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

func GetUserPointDetailList(ctx context.Context, filter model.UserPointDetailFilter) (int64, []model.UserPointDetailItem, error) {
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	opts := buildUserPointDetailDBOptionFromFilter(userPointDetailQuery, filter)
	userPointDetailPOs, err := userPointDetailQuery.GetUserPointDetail(ctx, opts...)
	if err != nil {
		return 0, nil, err
	}

	totalValue, err := userPointDetailQuery.GetUserPoint(ctx, filter.UserID)
	if err != nil {
		return 0, nil, err
	}

	result := make([]model.UserPointDetailItem, 0)
	for _, detailPO := range userPointDetailPOs {
		result = append(result, converter.ConvertUserPointDetailItemFromPO(detailPO))
	}
	return totalValue, result, nil
}

func GetUserPointDetailCount(ctx context.Context, filter model.UserPointDetailFilter) (int64, error) {
	userPointDetailQuery := repository.NewUserPointDetailQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildUserPointDetailDBOptionFromFilter(userPointDetailQuery, filter)
	return userPointDetailQuery.GetUserPointDetailCount(ctx, opts...)
}

// HINT: 以下的几个UserPoint相关函数都是并发安全的, 但不保证成功，事务失败时需要上层自行处理
func ChangeUserPoints(ctx context.Context, userID int64, eventType types.PointEventType, value int64, description string) error {
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
		err := userQuery.UpdateUser(ctx, user, repository.WithOptimisticLock("points", originalPoints))
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, userID, eventType, value, description)
		if err != nil {
			return err
		}
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
		if user.ID == (senderID) {
			senderPO = &user
		}
		if user.ID == (receiverID) {
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
	operations := func(repo repository.IRepository) error {
		userQuery := repo.NewUserQuery()
		userPointDetailQuery := repo.NewUserPointQuery()
		err := userQuery.UpdateUser(ctx, *senderPO, repository.WithOptimisticLock("points", senderOriginalPoints))
		if err != nil {
			return err
		}
		err = userQuery.UpdateUser(ctx, *receiverPO, repository.WithOptimisticLock("points", receiverOriginalPoints))
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, senderID, types.PointEventTransfer, -value, description)
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, receiverID, types.PointEventTransfer, receivedValue, description)
		if err != nil {
			return err
		}
		return nil
	}
	err = repo.InTransaction(ctx, operations)
	return err
}
