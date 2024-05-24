package leaderboard

import "gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"

type Leaderboard struct {
	Drivers []driver.Driver `json:"drivers"`
}
