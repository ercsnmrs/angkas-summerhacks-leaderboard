package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/trip"
)

type tripWriter interface {
	UpdateLeaderboard(ctx context.Context, trip trip.Event) error
}

func ConsumeTripCompleted(w tripWriter) JobHandler {
	return func(ctx context.Context, job Job) error {
		var d trip.Event
		if err := json.Unmarshal(job.Payload, &d); err != nil {
			return fmt.Errorf("json unmarshall: %s", err)
		}

		if d.Status == "complete" {
			if err := w.UpdateLeaderboard(ctx, d); err != nil {
				return fmt.Errorf("failed to update leaderboard: %s", err)
			}
		}

		return nil
	}
}
