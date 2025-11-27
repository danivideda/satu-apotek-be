package job

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/go-co-op/gocron/v2"
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

func (s *MyScheduler) AddClearCacheJob() {
	job, err := s.scheduler.NewJob(gocron.DurationJob(15*time.Second), gocron.NewTask(func(ctx context.Context) {
		items, err := s.repo.ApotekCode.DeleteExpired(ctx)
		if err != nil {
			s.logger.Error(err.Error())
		}

		if len(*items) > 0 {
			itemIDs := make([]string, len(*items))
			for idx, item := range *items {
				itemIDs[idx] = item.ApotekID.String()
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

func (s *MyScheduler) Start() {
	s.scheduler.Start()
}

func (s *MyScheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
