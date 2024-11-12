package repository

import (
	"context"
	"fmt"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"log"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func InitTestDB(t *testing.T) (IRepository, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	repo := NewRepository(db)
	err = Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	return repo, db
}
func ClearTestDB(db *gorm.DB) error {
	var tables []string
	// Retrieve the list of tables
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name != 'sqlite_sequence'").Scan(&tables).Error; err != nil {
		return err
	}
	// Drop each table
	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			return err
		}
	}
	return nil
}
func InitTestUser(t *testing.T, repo IRepository, n int) ([]po.UserPO, error) {
	ctx := context.Background()
	userQuery := repo.NewUserQuery()
	userPOs := make([]po.UserPO, 0)
	for i := 0; i < n; i++ {
		email := fmt.Sprintf("test:transaction:%d@example.com", i)
		password := fmt.Sprintf("test:transaction:%d", i)
		users, err := userQuery.GetUser(ctx, WithEmail(email))
		var user *po.UserPO
		if err == nil && len(users) != 0 {
			user = &users[0]
			user.Points = 1000
			err = userQuery.UpdateUser(ctx, *user)
			if err != nil {
				return nil, err
			}
		} else {
			user, err = userQuery.CreateUser(ctx, email, password)
			if user == nil || err != nil {
				return nil, err
			}
			user.Points = 1000
			err := userQuery.UpdateUser(ctx, *user)
			if err != nil {
				return nil, err
			}

		}
		userPOs = append(userPOs, *user)
	}
	for _, user := range userPOs {
		assert.Equal(t, int64(1000), user.Points)
	}
	return userPOs, nil
}
func ResetUserPoints(users []po.UserPO) {
	for _, users := range users {
		users.Points = 1000
	}
}
func QueryUser(repo IRepository, idx int) (po.UserPO, []po.UserPointDetailPO, error) {
	userQuery := repo.NewUserQuery()
	userPointDetailQuery := repo.NewUserPointQuery()
	ctx := context.Background()
	email := fmt.Sprintf("test:transaction:%d@example.com", idx)
	userPOs, err := userQuery.GetUser(ctx, WithEmail(email))
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

func TestInTransAction(t *testing.T) {
	ctx := context.Background()

	TransferOpsSucceed := func(repo IRepository) error {
		userPOs, err := InitTestUser(t, repo, 2)
		assert.Nil(t, err)
		assert.Len(t, userPOs, 2)
		user1 := userPOs[0]
		user2 := userPOs[1]
		userQuery := repo.NewUserQuery()
		userPointDetailQuery := repo.NewUserPointQuery()
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		err = userQuery.UpdateUser(ctx, user1)
		if err != nil {
			return err
		}
		err = userQuery.UpdateUser(ctx, user2)
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), model.PointEventTransfer, -100, "test")
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), model.PointEventTransfer, 99, "test")
		if err != nil {
			return err
		}
		return nil
	}
	TransferOpsRollback1 := func(repo IRepository) error {
		userPOs, err := InitTestUser(t, repo, 2)
		assert.Nil(t, err)
		assert.Len(t, userPOs, 2)
		user1 := userPOs[0]
		user2 := userPOs[1]
		userQuery := repo.NewUserQuery()
		userPointDetailQuery := repo.NewUserPointQuery()
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		err = userQuery.UpdateUser(ctx, user1)
		if err != nil {
			return err
		}
		err = userQuery.UpdateUser(ctx, user2)
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), model.PointEventTransfer, -100, "test")
		if err != nil {
			return err
		}
		err = userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), model.PointEventTransfer, 99, "test")
		if err != nil {
			return err
		}
		return errors.New("test rollback at end")
	}
	TransferOpsRollback2 := func(repo IRepository) error {
		userPOs, err := InitTestUser(t, repo, 2)
		if err != nil {
			log.Printf("InitTestUser error: %v", err)
		}
		assert.Nil(t, err)
		assert.Len(t, userPOs, 2)
		user1 := userPOs[0]
		user2 := userPOs[1]
		userQuery := repo.NewUserQuery()
		userPointDetailQuery := repo.NewUserPointQuery()
		user1.Points -= 100
		t.Logf("user1 points: %d\n", user1.Points)
		user2.Points += 99
		t.Logf("user2 points: %d\n", user2.Points)
		err = userQuery.UpdateUser(ctx, user1)
		if err != nil {
			return err
		}
		err = userQuery.UpdateUser(ctx, user2)
		if err != nil {
			return err
		}
		err = errors.New("test rollback in mid")
		ign := userPointDetailQuery.CreateUserPointDetail(ctx, int64(user1.ID), model.PointEventTransfer, -100, "test")
		if ign != nil {
			return err
		}
		ign = userPointDetailQuery.CreateUserPointDetail(ctx, int64(user2.ID), model.PointEventTransfer, 99, "test")
		if ign != nil {
			return err
		}
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
			repo, db := InitTestDB(t)
			defer func() {
				err := ClearTestDB(db)
				if err != nil {
					t.Fatal(err)
				}
			}()

			terr := repo.InTransaction(ctx, tt.ops)
			tt.wantErr(t, terr, fmt.Sprintf("InTransAction(%v)", tt.ops))

			if terr == nil {
				// succeed
				user1PO, userPointDetails, err := QueryUser(repo, 0)
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
				assert.Equal(t, model.PointEventTransfer, userPointDetails[0].EventType)
				assert.Equal(t, int64(-100), userPointDetails[0].Value)
				user2PO, userPointDetails, err := QueryUser(repo, 1)
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
				assert.Equal(t, model.PointEventTransfer, userPointDetails[0].EventType)
				assert.Equal(t, int64(99), userPointDetails[0].Value)
			} else {
				// rollback
				// rollback create user
				_, _, err := QueryUser(repo, 0)
				assert.Error(t, err)
				_, _, err = QueryUser(repo, 1)
				assert.Error(t, err)
			}
		})
	}
}
