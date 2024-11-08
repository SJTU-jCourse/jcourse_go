package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"jcourse_go/dal"
	"jcourse_go/model/po"
	"log"
	"sync"
	"testing"
	"time"
)

func TestUserQuery_UpdateUser(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewUserQuery(db)
	t.Run("none", func(t *testing.T) {
		err := query.UpdateUser(ctx, po.UserPO{}, WithID(-1))
		log.Printf("err: %v", err)
		assert.Errorf(t, err, "error: %v", err)
	})

	t.Run("concurrently_opt_lock", func(t *testing.T) {
		var wg sync.WaitGroup
		const numRoutines = 100
		testUserEmail := "update_concurrent_test@example.com"
		users, err := query.GetUser(ctx, WithEmail(testUserEmail))
		assert.Nil(t, err)
		var user po.UserPO
		if len(users) != 0 {
			user = users[0]
		} else {
			pUser, err := query.CreateUser(ctx, testUserEmail, "password")
			user = *pUser
			assert.Nil(t, err)
		}

		if err != nil {
			return
		}

		originalPoints := user.Points
		errChan := make(chan error, numRoutines)
		for range numRoutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(100 * time.Millisecond)
				newUser := user
				newUser.Points += 100
				conErr := query.UpdateUser(ctx, newUser, WithOptimisticLock("points", originalPoints))
				if conErr != nil {
					errChan <- conErr
				}
			}()
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for conErr := range errChan {
			log.Printf("err: %v", conErr)
			errCount += 1
		}
		log.Printf("%v routine failed to update", errCount)
		assert.Equal(t, 99, errCount)
		users, err = query.GetUser(ctx, WithEmail(testUserEmail))
		assert.Nil(t, err)
		comp := func() bool { return len(users) > 0 }
		assert.Condition(t, comp)
		assert.Equal(t, users[0].Points, originalPoints+100)
	})
	t.Run("concurrently_no_lock", func(t *testing.T) {
		var wg sync.WaitGroup
		const numRoutines = 100
		testUserEmail := "update_concurrent_test@example.com"
		users, err := query.GetUser(ctx, WithEmail(testUserEmail))
		assert.Nil(t, err)
		var user po.UserPO
		if len(users) != 0 {
			user = users[0]
		} else {
			pUser, err := query.CreateUser(ctx, testUserEmail, "password")
			user = *pUser
			assert.Nil(t, err)
		}

		if err != nil {
			return
		}
		errChan := make(chan error, numRoutines)
		for i := range numRoutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(100 * time.Millisecond)
				newUser := user
				newUser.Points += int64(i)
				conErr := query.UpdateUser(ctx, newUser)
				if conErr != nil {
					errChan <- conErr
				}
			}()
		}

		wg.Wait()
		close(errChan)
		errCount := 0
		for conErr := range errChan {
			log.Printf("err: %v", conErr)
			errCount += 1
		}
		log.Printf("%v routine failed to update", errCount)
		assert.Equal(t, 0, errCount)
		users, err = query.GetUser(ctx, WithEmail(testUserEmail))
		assert.Nil(t, err)
		log.Printf("points: %d", users[0].Points)
		// 大概率不是100
	})
}
