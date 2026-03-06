package job

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type MyScheduler struct {
	scheduler gocron.Scheduler
	logger    gocron.Logger
	repo      repository.Repository
}

func NewScheduler(r repository.Repository) (*MyScheduler, error) {
	l := gocron.NewLogger(gocron.LogLevelInfo)
	s, err := gocron.NewScheduler(gocron.WithLogger(l))
	if err != nil {
		return nil, err
	}

	scheduler := MyScheduler{
		scheduler: s,
		logger:    l,
		repo:      r,
	}
	return &scheduler, nil
}

func (s *MyScheduler) AddClearApotekCodeJob() {
	job, err := s.scheduler.NewJob(gocron.DurationJob(5*time.Second), gocron.NewTask(func(ctx context.Context) {
		items, err := s.repo.ApotekCode.DeleteExpired(ctx)
		if err != nil {
			s.logger.Error(err.Error())
		}

		if len(*items) > 0 {
			itemIDs := make([]int64, len(*items))
			for idx, item := range *items {
				itemIDs[idx] = item.ApotekID
			}
			s.logger.Info("job: deleted expired Apotek Code row(s)", itemIDs)
		} else {
			s.logger.Info("job: no expired Apotek Code")
		}
	}))
	if err != nil {
		s.logger.Error(err.Error())
	}
	s.logger.Info("gocron: job added", job.ID())
}

func (s *MyScheduler) AddDeleteExpiredSessionsJob() {
	jobDuration := gocron.DurationJob(5 * time.Second)
	task := gocron.NewTask(func(ctx context.Context) {
		items, err := s.repo.OwnerSessions.DeleteExpired(ctx)
		if err != nil {
			s.logger.Error(err.Error())
		}

		if len(*items) > 0 {
			ownerSessionIDs := make([]pgtype.UUID, len(*items))
			for idx, item := range *items {
				ownerSessionIDs[idx] = item.ID
			}
			s.logger.Info("job: deleted expired Owner Session row(s)", ownerSessionIDs)
		} else {
			s.logger.Info("job: no expired owner session")
		}
	})
	job, err := s.scheduler.NewJob(jobDuration, task)
	if err != nil {
		s.logger.Error(err.Error())
	}
	s.logger.Info("delete expired owner session job added", job.ID())
}

func (s *MyScheduler) Start() {
	s.scheduler.Start()
}

func (s *MyScheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
