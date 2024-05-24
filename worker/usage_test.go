package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func Test_usageLimiter_Use(t *testing.T) {
	tests := []struct {
		name string
		// deps
		source sourceLimiter
		// params
		maxUsage int
		// returns
		wantRemaining int
		wantErr       bool
	}{
		{
			"initial usage",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 0, nil
				},
				UpdateFn: func(ctx context.Context, key string, val int, expr time.Duration) error {
					return nil
				},
			},
			3,
			2,
			false,
		},
		{
			"2nd usage",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 1, nil
				},
				UpdateFn: func(ctx context.Context, key string, val int, expr time.Duration) error {
					return nil
				},
			},
			3,
			1,
			false,
		},
		{
			"last usage",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 2, nil
				},
				UpdateFn: func(ctx context.Context, key string, val int, expr time.Duration) error {
					return nil
				},
			},
			3,
			0,
			false,
		},
		{
			"exhausted",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 3, nil
				},
			},
			3,
			0,
			true,
		},
		{
			"overloaded",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 4, nil
				},
			},
			3,
			0,
			true,
		},
		{
			"failed count check",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 0, errors.New("source failure")
				},
			},
			3,
			0,
			true,
		},
		{
			"failed usage update",
			&mockLimiterSource{
				CountFn: func(ctx context.Context, key string) (int, error) {
					return 1, nil
				},
				UpdateFn: func(ctx context.Context, key string, val int, expr time.Duration) error {
					return errors.New("source failure")
				},
			},
			3,
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &usageLimiter{
				prefix:   "test",
				maxUsage: tt.maxUsage,
				source:   tt.source,
			}
			ctx := context.Background()
			gotRemaining, err := l.Use(ctx, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("Use() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRemaining != tt.wantRemaining {
				t.Errorf("Use() gotRemaining = %v, want %v", gotRemaining, tt.wantRemaining)
			}
		})
	}
}

func Test_ComputeResetOffset(t *testing.T) {
	loc := time.Now().Location()
	tests := []struct {
		name       string
		t          time.Time
		resetHour  string
		wantOffset time.Duration
	}{
		{
			"5:30 AM - 3:00 AM",
			time.Date(2023, 9, 14, 5, 30, 0, 0, loc),
			"3:00AM",
			21*time.Hour + 30*time.Minute,
		},
		{
			"2:00 PM - 3:00 AM",
			time.Date(2023, 9, 14, 14, 0, 0, 0, loc),
			"3:00AM",
			13 * time.Hour,
		},
		{
			"11:00 PM - 3:00 AM",
			time.Date(2023, 9, 14, 23, 0, 0, 0, loc),
			"3:00AM",
			4 * time.Hour,
		},
		{
			"1:00 AM - 3:00 AM",
			time.Date(2023, 9, 14, 1, 0, 0, 0, loc),
			"3:00AM",
			2 * time.Hour,
		},
		{
			"12:00 AM - 3:00 AM",
			time.Date(2023, 9, 14, 0, 0, 0, 0, loc),
			"3:00AM",
			3 * time.Hour,
		},
		{
			"4:00 AM - 3:00 AM",
			time.Date(2023, 9, 14, 4, 0, 0, 0, loc),
			"3:00AM",
			23 * time.Hour,
		},
		{
			"3:00 AM - 3:00 AM",
			time.Date(2023, 9, 14, 3, 0, 0, 0, loc),
			"3:00AM",
			24 * time.Hour,
		},
		{
			"2:59 AM - 3:00 AM",
			time.Date(2023, 9, 14, 2, 59, 0, 0, loc),
			"3:00AM",
			time.Minute,
		},
		{
			"3:01 AM - 3:00 AM",
			time.Date(2023, 9, 14, 3, 1, 0, 0, loc),
			"3:00AM",
			24*time.Hour - time.Minute,
		},
		{
			"3:01 AM - 3:30 AM",
			time.Date(2023, 9, 14, 3, 1, 0, 0, loc),
			"3:30AM",
			29 * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset, err := time.Parse(time.Kitchen, tt.resetHour)
			if err != nil {
				t.Errorf("time error format: %s", err)
				return
			}
			if got := ComputeResetOffset(tt.t, reset); got != tt.wantOffset {
				t.Errorf("ComputeResetOffset got: %s, want: %s", got, tt.wantOffset)
			}
		})
	}
}

type mockLimiterSource struct {
	CountFn  func(ctx context.Context, key string) (int, error)
	UpdateFn func(ctx context.Context, key string, val int, expr time.Duration) error
}

func (m *mockLimiterSource) Count(ctx context.Context, key string) (int, error) {
	return m.CountFn(ctx, key)
}

func (m *mockLimiterSource) Update(ctx context.Context, key string, val int, expr time.Duration) error {
	return m.UpdateFn(ctx, key, val, expr)
}
