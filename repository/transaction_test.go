package repository

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"jcourse_go/dal"
	"jcourse_go/model/po"
	"log"
	"testing"
)

//	func NewTestDB() *gorm.DB {
//		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
//		if err != nil {
//			panic("failed to connect database")
//		}
//		return db
//	}
func InitTestUser(db *gorm.DB, n int) []po.UserPO {
	ctx := context.Background()
	userdb := db.WithContext(ctx).Find(&po.UserPO{})
	for i := 0; i < n; i++ {
		userdb.Create(&po.UserPO{
			Username: fmt.Sprintf("test:transaction:%d", i),
			Password: fmt.Sprintf("test:transaction:%d", i),
			Email:    fmt.Sprintf("test:transaction:%d@example.com", i),
			UserRole: "student",
			Points:   1000,
		})
	}
	userQuery := NewUserQuery(db)
	userPOs := make([]po.UserPO, 0)
	for i := 0; i < n; i++ {
		userPO, err := userQuery.GetUser(ctx, WithEmail(fmt.Sprintf("test:transaction:%d@example.com", i)))
		if err != nil {
			fmt.Printf("error: %v\n", err)
			panic(err)
		}

		if len(userPO) == 0 {
			_ = fmt.Errorf("user not found")
			panic("user not found")
		}
		userPOs = append(userPOs, userPO[0])
	}
	return userPOs
}
func QueryUser(db *gorm.DB, idx int) (po.UserPO, []po.UserPointDetailPO, error) {
	userQuery := NewUserQuery(db)
	userPointDetailQuery := NewUserPointDetailQuery(db)
	ctx := context.Background()
	userPOs, err := userQuery.GetUser(ctx, WithEmail(fmt.Sprintf("test:transaction:%d@example.com", idx)))
	if err != nil {
		return po.UserPO{}, nil, err
	}
	if len(userPOs) == 0 {
		return po.UserPO{}, nil, errors.Errorf("user not found")
	}
	userPO := userPOs[0]
	userPointDetails, err := userPointDetailQuery.GetUserPointDetail(ctx, WithUserID(int64(userPO.ID)))
	if err != nil {
		return userPO, nil, err
	}
	return userPO, userPointDetails, nil
}

func CleanupTestUser(db *gorm.DB, n int) {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		log.Fatalf("Failed to start transaction: %v", tx.Error)
	}

	for i := 0; i < n; i++ {
		if err := tx.Where("email = ?", fmt.Sprintf("test:transaction:%d@example.com", i)).Delete(&po.UserPO{}).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to delete user: %v", err)
		}
		if err := tx.Where("description = ?", "test").Delete(&po.UserPointDetailPO{}).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to delete user point detail: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}
func TestInTransAction(t *testing.T) {
	db := dal.GetDBClient()
	handler := NewTransactionHandler(db)
	ctx := context.Background()
	err := Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	users := InitTestUser(db, 2)
	user1 := users[0]
	assert.Equal(t, user1.Points, int64(1000))
	user2 := users[1]
	assert.Equal(t, user2.Points, int64(1000))
	defer CleanupTestUser(db, 2)
	TransferOpsSucceed := func(db *gorm.DB) error {
		userQuery := NewUserQuery(db)
		userPointDetailQuery := NewUserPointDetailQuery(db)
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		userQuery.UpdateUser(ctx, user1)
		userQuery.UpdateUser(ctx, user2)
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), po.PointEventTransfer, -100, "test")
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), po.PointEventTransfer, 99, "test")
		return nil
	}
	TransferOpsRollback1 := func(db *gorm.DB) error {
		userQuery := NewUserQuery(db)
		userPointDetailQuery := NewUserPointDetailQuery(db)
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		userQuery.UpdateUser(ctx, user1)
		userQuery.UpdateUser(ctx, user2)
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), po.PointEventTransfer, -100, "test")
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), po.PointEventTransfer, 99, "test")
		return errors.Errorf("test rollback at end")
	}
	TransferOpsRollback2 := func(db *gorm.DB) error {
		userQuery := NewUserQuery(db)
		userPointDetailQuery := NewUserPointDetailQuery(db)
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		userQuery.UpdateUser(ctx, user1)
		userQuery.UpdateUser(ctx, user2)
		err = errors.Errorf("test rollback in mid")
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), po.PointEventTransfer, -100, "test")
		userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), po.PointEventTransfer, 99, "test")
		return err
	}
	tests := []struct {
		name    string
		ops     DBOperation
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{name: "transfer", ops: TransferOpsSucceed, wantErr: assert.NoError},
		{name: "rollback1", ops: TransferOpsRollback1, wantErr: assert.Error},
		{name: "rollback2", ops: TransferOpsRollback2, wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terr := InTransAction(ctx, handler, tt.ops)
			tt.wantErr(t, terr, fmt.Sprintf("InTransAction(%v)", tt.ops))

			if terr == nil {
				// succeed
				user1PO, userPointDetails, err := QueryUser(db, 0)
				if err != nil {
					t.Fatal(err)
					return
				}
				if len(userPointDetails) == 0 {
					t.Fatal("user1 point details not found")
					return
				}
				assert.Equal(t, int64(900), user1PO.Points)
				assert.Len(t, userPointDetails, 1)
				assert.Equal(t, int64(user1PO.ID), userPointDetails[0].UserID)
				assert.Equal(t, po.PointEventTransfer, userPointDetails[0].EventType)
				assert.Equal(t, int64(-100), userPointDetails[0].Value)
				user2PO, userPointDetails, err := QueryUser(db, 1)
				if err != nil {
					t.Fatal(err)
					return
				}
				if len(userPointDetails) == 0 {
					t.Fatal("user1 point details not found")
					return
				}
				assert.Equal(t, int64(1099), user2PO.Points)
				assert.Len(t, userPointDetails, 1)
				assert.Equal(t, int64(user2PO.ID), userPointDetails[0].UserID)
				assert.Equal(t, po.PointEventTransfer, userPointDetails[0].EventType)
				assert.Equal(t, int64(99), userPointDetails[0].Value)
			} else {
				// rollback
				user1PO, userPointDetails, err := QueryUser(db, 0)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, int64(1000), user1PO.Points)
				assert.Len(t, userPointDetails, 0)
				user2PO, userPointDetails, err := QueryUser(db, 1)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, int64(1000), user2PO.Points)
				assert.Len(t, userPointDetails, 0)
			}

		})
	}
}
