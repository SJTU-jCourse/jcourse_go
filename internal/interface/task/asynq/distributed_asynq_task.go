package asynq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/XQ-Gang/gg/gptr"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"jcourse_go/internal/interface/task/base"
	"jcourse_go/internal/interface/task/lock"
)

type distributedScheduler struct {
	location   *time.Location
	cronEngine *cron.Cron

	mu    sync.Mutex
	idMap map[string]cron.EntryID

	lockMgr lock.IDistributedLockManager
	tm      *AsynqTaskManager
	killSub *redis.PubSub
}

type distributedJob struct {
	id         string
	task       base.Task
	interval   time.Duration
	lockMgr    lock.IDistributedLockManager
	tm         *AsynqTaskManager
	lockMargin time.Duration
}

func (j *distributedJob) Run() {
	lockKey := "dist_lock:" + j.task.Type()
	ttl := j.interval - j.lockMargin // release the lock in advanced
	if ttl < time.Second {
		ttl = time.Second
	}

	dLock, err := j.lockMgr.NewLock(lockKey, ttl)
	if err != nil {
		log.Printf("NewLock Failed! error: %+v \n", err)
		return
	}

	acquired, err := dLock.TryLock()
	if err != nil {
		log.Printf("TryLock Failed! error: %+v \n", err)
		return
	}
	if !acquired {
		// abort the task if other clients have already acquired the lock
		log.Printf("Lock already acquired by other client, aborting task %s \n", j.task.Type())
		return
	}

	err = j.tm.Enqueue(j.task, base.TaskOption{
		WithQueue: gptr.Of(DistributedPeriodicQueue), // specfic queue
	})
	if err != nil {
		log.Printf("Enqueue Failed! error: %+v \n", err)
		_, _ = dLock.Unlock()
	} else {
		fmt.Println("Enqueue Success! Current Time: ", time.Now().Format(time.RFC3339))
	}
}

func newDistributedScheduler(lockMgr lock.IDistributedLockManager, tm *AsynqTaskManager, rdb *redis.Client) *distributedScheduler {
	loc := time.UTC
	s := &distributedScheduler{
		location:   loc,
		cronEngine: cron.New(cron.WithLocation(loc)),
		idMap:      make(map[string]cron.EntryID),
		lockMgr:    lockMgr,
		tm:         tm,
		killSub:    rdb.Subscribe(context.Background(), "task_kill_channel"),
	}

	go s.listenForKills()

	return s
}

func (s *distributedScheduler) listenForKills() {
	for {
		msg, err := s.killSub.ReceiveMessage(context.Background())
		if err != nil {
			if errors.Is(err, redis.ErrClosed) {
				log.Println("killSub closed, stopping listenForKills")
				break
			}
			log.Printf("Failed to receive message: %+v \n", err)
			continue
		}
		taskID := msg.Payload
		err = s.unregister(taskID)
		if err != nil {
			log.Printf("Failed to unregister task: %+v \n", err)
		}
	}
}

func (s *distributedScheduler) startScheduler() {
	s.cronEngine.Start()
}

func (s *distributedScheduler) shutdown() error {
	s.cronEngine.Stop()
	err := s.killSub.Close()
	if err != nil {
		log.Printf("Failed to close redis client: %+v \n", err)
		return err
	}
	return nil
}

func (s *distributedScheduler) register(interval base.TaskInterval, job *distributedJob) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.idMap[job.id]; ok {
		return "", fmt.Errorf("asynq: job ID %s already exists", job.id)
	}

	cronSpec := fmt.Sprintf("@every %s", interval.String())
	log.Printf("Adding job %s with interval %s \n", job.id, cronSpec)

	entryID, err := s.cronEngine.AddJob(cronSpec, job)
	if err != nil {
		return "", err
	}
	s.idMap[job.id] = entryID
	return job.id, nil
}

func (s *distributedScheduler) unregister(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cronID, ok := s.idMap[jobID]
	if !ok {
		return fmt.Errorf("asynq: no scheduler entry found for job ID %s", jobID)
	}
	s.cronEngine.Remove(cronID)
	delete(s.idMap, jobID)
	return nil
}

type DistributedAsynqTaskManager struct {
	*AsynqTaskManager
	lock.IDistributedLockManager
	*distributedScheduler

	rdb *redis.Client
}

func NewDistributedAsynqTaskManager(
	tm *AsynqTaskManager,
	rdb *redis.Client,
) *DistributedAsynqTaskManager {
	lockMgr := lock.NewRedisDistributedLockManager(rdb)
	m := &DistributedAsynqTaskManager{
		AsynqTaskManager:        tm,
		IDistributedLockManager: lockMgr,
		rdb:                     rdb,
		distributedScheduler: newDistributedScheduler(
			lockMgr,
			tm,
			rdb,
		),
	}
	m.distributedScheduler.startScheduler()
	return m
}

func (m *DistributedAsynqTaskManager) Submit(interval base.TaskInterval, task base.Task) (base.PeriodicTaskId, error) {
	if m.scheduler == nil {
		return "", errors.New("failed to initialize distributed scheduler")
	}
	if interval <= 0 {
		return "", errors.New("interval must be positive")
	}

	jobID := fmt.Sprintf("%s_%s", "job_id", task.Type()) // keep consistent among all clients
	job := &distributedJob{
		id:         jobID,
		task:       task,
		interval:   interval,
		lockMgr:    m.IDistributedLockManager,
		tm:         m.AsynqTaskManager,
		lockMargin: 10 * time.Second, // optional margin for TTL
	}

	_, err := m.distributedScheduler.register(interval, job)
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (m *DistributedAsynqTaskManager) Kill(taskID base.PeriodicTaskId) error {
	err := m.rdb.Publish(context.Background(), "task_kill_channel", taskID).Err()
	if err != nil {
		return fmt.Errorf("failed to publish kill notification for task %s: %v", taskID, err)
	}

	return nil
}

func (m *DistributedAsynqTaskManager) Shutdown() error {
	if m.distributedScheduler != nil {
		err := m.distributedScheduler.shutdown()
		if err != nil {
			return err
		}
	}
	return m.AsynqTaskManager.Shutdown()
}
