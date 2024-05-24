package worker

import (
	"context"
	"fmt"
	"time"
)

type usageLimiter struct {
	prefix     string
	resetClock time.Time
	maxUsage   int
	source     sourceLimiter
}

func (l *usageLimiter) Use(ctx context.Context, id string) (remaining int, err error) {
	key := fmt.Sprintf("%s:%s", l.prefix, id)

	uses, err := l.source.Count(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("source.Uses: %s", err)
	}

	if uses >= l.maxUsage {
		return 0, fmt.Errorf("max usage reached %d/%d", uses, l.maxUsage)
	}

	uses += 1
	remaining = l.maxUsage - uses
	n := time.Now()
	expr := ComputeResetOffset(n, l.resetClock)
	if err = l.source.Update(ctx, key, uses, expr); err != nil {
		return 0, fmt.Errorf("source.Update: %s", err)
	}

	return remaining, nil
}

func (l *usageLimiter) Remaining(ctx context.Context, id string) (remaining int, err error) {
	key := fmt.Sprintf("%s:%s", l.prefix, id)

	uses, err := l.source.Count(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("source.Uses: %s", err)
	}
	return l.maxUsage - uses, nil
}

// func newLimitReset(name string, resetClock time.Time, maxUsageDaily int, src sourceLimiter) *usageLimiter {
// 	return &usageLimiter{name, resetClock, maxUsageDaily, src}
// }

type sourceLimiter interface {
	Count(ctx context.Context, key string) (int, error)
	Update(ctx context.Context, key string, val int, expr time.Duration) error
}

// ComputeResetOffset returns duration offset on resetHour the next day.
// ex. given 3:00 as resetHour with current time 2023-09-14 05:30 should reset the next day
// at 2023-09-15 03:00 with 21h30m as offset.
func ComputeResetOffset(t time.Time, reset time.Time) time.Duration {
	next := time.Date(t.Year(), t.Month(), t.Day()+1, reset.Hour(), reset.Minute(), 0, 0, t.Location())
	diff := next.Sub(t)

	// Catch edge case that it passed midnight and need to reset at resetHour.
	rd := reset.Hour()*60 + reset.Minute()
	td := t.Hour()*60 + t.Minute()
	if td < rd {
		return time.Duration(rd-td) * time.Minute
	}
	return diff
}
