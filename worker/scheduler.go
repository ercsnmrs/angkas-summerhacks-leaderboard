package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Schedule struct {
	Reset     time.Time
	Frequency time.Duration
	Fn        func(ctx context.Context) error

	done   chan struct{}
	logger *slog.Logger
}

func (s *Schedule) run(ctx context.Context) error {
	starts := ComputeResetOffset(time.Now(), s.Reset)
	s.logger.Info(fmt.Sprintln("schedule starts in", starts))
	t := time.NewTimer(starts)
	<-t.C

	s.logger.Info(fmt.Sprintln("schedule started at", t))
	ticker := time.NewTicker(s.Frequency)
	go func() {
		for {
			select {
			case <-s.done:
				s.logger.Info("schedule exited")
				return
			case tk := <-ticker.C:
				s.logger.Info(fmt.Sprintln("schedule executed at", tk))
				if err := s.Fn(ctx); err != nil {
					s.logger.Error("schedule error", "err", err)
				}
			}
		}
	}()

	return nil
}

func NewSchedule(clock string, frequency time.Duration, fn func(ctx context.Context) error) (Schedule, error) {
	reset := time.Now().Add(7 * time.Second)

	if clock != "" {
		parsedReset, err := time.Parse(time.Kitchen, clock)
		if err != nil {
			return Schedule{}, fmt.Errorf("error parsing clock: %w", err)
		}
		reset = parsedReset
	}

	return Schedule{
		Reset:     reset,
		Frequency: frequency,
		Fn:        fn,
	}, nil
}