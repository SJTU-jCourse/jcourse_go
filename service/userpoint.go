package service

import (
	"context"
	"fmt"
	"time"

	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
	"jcourse_go/repository"
	"jcourse_go/util"

	"github.com/pkg/errors"
)

func buildUserPointDetailDBOptionFromFilter(ctx context.Context, q *repository.Query, filter model.UserPointDetailFilter) repository.IUserPointDetailPODo {
	builder := q.UserPointDetailPO.WithContext(ctx)
	p := q.UserPointDetailPO

	if filter.UserPointDetailID > 0 {
		builder = builder.Where(p.ID.Eq(filter.UserPointDetailID))
	}
	if filter.UserID > 0 {
		builder = builder.Where(p.UserID.Eq(filter.UserID))
	}
	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.StartTime > 0 && filter.EndTime > 0 {
		builder = builder.Where(p.CreatedAt.Between(time.Unix(filter.StartTime, 0), time.Unix(filter.EndTime, 0)))
	} else if filter.StartTime > 0 {
		builder = builder.Where(p.CreatedAt.Gte(time.Unix(filter.StartTime, 0)))
	} else if filter.EndTime > 0 {
		builder = builder.Where(p.CreatedAt.Lte(time.Unix(filter.EndTime, 0)))
	}
	return builder
}

func GetUserPointDetailList(ctx context.Context, filter model.UserPointDetailFilter) (int64, []model.UserPointDetailItem, error) {

	p := repository.Q.UserPointDetailPO
	q := buildUserPointDetailDBOptionFromFilter(ctx, repository.Q, filter)
	userPointDetailPOs, err := q.Find()
	if err != nil {
		return 0, nil, err
	}

	total := struct {
		Value int64 `json:"value"`
	}{}
	err = repository.Q.UserPointDetailPO.WithContext(ctx).Select(p.Value.Sum().As("total")).Where(p.UserID.Eq(filter.UserID)).Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	result := make([]model.UserPointDetailItem, 0)
	for _, detailPO := range userPointDetailPOs {
		result = append(result, converter.ConvertUserPointDetailItemFromPO(*detailPO))
	}
	return total.Value, result, nil
}

func GetUserPointDetailCount(ctx context.Context, filter model.UserPointDetailFilter) (int64, error) {
	filter.Page, filter.PageSize = 0, 0
	q := buildUserPointDetailDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

// HINT: 以下的几个UserPoint相关函数都是并发安全的, 但不保证成功，事务失败时需要上层自行处理
func ChangeUserPoints(ctx context.Context, userID int64, eventType types.PointEventType, value int64, description string) error {
	u := repository.Q.UserPO
	user, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).Take()
	if err != nil {
		return err
	}

	if user.Points+value < 0 {
		return errors.Errorf("user %d has not enough points", userID)
	}
	originalPoints := user.Points
	user.Points += value

	point := po.UserPointDetailPO{
		UserID:      userID,
		Value:       value,
		Description: description,
		EventType:   string(eventType),
	}

	err = repository.Q.Transaction(func(tx *repository.Query) error {
		_, err := tx.UserPO.WithContext(ctx).Where(u.ID.Eq(user.ID), u.Points.Eq(originalPoints)).Update(u.Points, user.Points)
		if err != nil {
			return err
		}
		err = tx.UserPointDetailPO.WithContext(ctx).Create(&point)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func calcHandlingFee(ctx context.Context, value int64) int64 {
	siteSetting, err := GetSetting(ctx, constant.HandleFeeRateKey)
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
	u := repository.Q.UserPO

	// 合并到一次查询
	ids := []int64{senderID, receiverID}
	userPOs, err := u.WithContext(ctx).Where(u.ID.In(ids...)).Find()
	if err != nil {
		return err
	}
	var senderPO *po.UserPO = nil
	var receiverPO *po.UserPO = nil
	for _, user := range userPOs {
		if user.ID == (senderID) {
			senderPO = user
		}
		if user.ID == (receiverID) {
			receiverPO = user
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

	description := fmt.Sprintf(TransferDescriptionFormat, senderID, receiverID, value)
	senderPoint := po.UserPointDetailPO{
		UserID:      senderID,
		Value:       value,
		Description: description,
		EventType:   string(types.PointEventTransfer),
	}
	receiverPoint := po.UserPointDetailPO{
		UserID:      receiverID,
		Value:       receivedValue,
		Description: description,
		EventType:   string(types.PointEventTransfer),
	}
	err = repository.Q.Transaction(func(tx *repository.Query) error {
		_, err := tx.UserPO.WithContext(ctx).Where(u.ID.Eq(senderID), u.Points.Eq(senderOriginalPoints)).Update(u.Points, senderPO.Points)
		if err != nil {
			return err
		}
		_, err = tx.UserPO.WithContext(ctx).Where(u.ID.Eq(receiverID), u.Points.Eq(receiverOriginalPoints)).Update(u.Points, receiverPO.Points)
		if err != nil {
			return err
		}

		err = tx.UserPointDetailPO.WithContext(ctx).Create(&senderPoint)
		if err != nil {
			return err
		}
		err = tx.UserPointDetailPO.WithContext(ctx).Create(&receiverPoint)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
