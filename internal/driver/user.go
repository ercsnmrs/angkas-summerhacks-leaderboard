package driver

import (
	"time"
)

type RFM struct {
	Recency   float64 `json:"recency"`
	Frequency float64 `json:"frequency"`
	Monetary  float64 `json:"monetary"`
}

type Rating struct {
	RFM     RFM     `json:"rating"`
	Average float64 `json:"average"`
}

type Driver struct {
	DriverID                     string    `json:"driver_id"`
	LastCompletedTripDate        time.Time `json:"last_completed_trip_date"`
	NetIncome                    float64   `json:"net_income"`
	NumberOfCompletedTrips       int       `json:"number_of_completed_trips"`
	UniqueDateWithCompletedTrips int       `json:"unique_date_with_completed_trips"`
	ServiceZone                  string    `json:"service_zone"`
	Rating                       Rating    `json:"rating"`
}

func (d Driver) CalculateAverage() float64 {
	// Cap the number of completed trips to a maximum of 5 points
	maxTrips := 5
	tripsScore := float64(d.NumberOfCompletedTrips)
	if tripsScore > float64(maxTrips) {
		tripsScore = float64(maxTrips)
	}
	normalizedTripsScore := (tripsScore / float64(maxTrips)) * 5

	// Normalize RFM values to be on a scale of 0 to 5
	normalizedRFM := ((d.Rating.RFM.Recency + d.Rating.RFM.Frequency + d.Rating.RFM.Monetary) / 3) * 5 / 4

	// Calculate the weighted average
	return (0.6 * normalizedTripsScore) + (0.4 * normalizedRFM)
}

// CalculateMonetary
func (d Driver) CalculateMonetary(currentHighestNetEarnings float64) float64 {
	highestNetEarnings := currentHighestNetEarnings
	pastMonthEarnings := d.NetIncome
	if pastMonthEarnings >= 0.75*highestNetEarnings {
		return 4
	} else if pastMonthEarnings >= 0.5*highestNetEarnings {
		return 3
	} else if pastMonthEarnings >= 0.25*highestNetEarnings {
		return 2
	} else if pastMonthEarnings >= 0.01*highestNetEarnings {
		return 1
	} else {
		return 0
	}
}

// Calculate Recency
func (d Driver) CalculateRecency() float64 {
	lastTripDate := d.LastCompletedTripDate
	daysSinceLastTrip := time.Since(lastTripDate).Hours() / 24

	if daysSinceLastTrip <= 7 {
		return 4
	} else if daysSinceLastTrip <= 14 {
		return 3
	} else if daysSinceLastTrip <= 21 {
		return 2
	} else if daysSinceLastTrip <= 30 {
		return 1
	} else {
		return 0
	}
}

// Calculate Frequency
func (d Driver) CalculateFrequency() float64 {
	activeDays := d.UniqueDateWithCompletedTrips
	if activeDays >= 22 && activeDays <= 30 {
		return 4
	} else if activeDays >= 15 && activeDays <= 21 {
		return 3
	} else if activeDays >= 8 && activeDays <= 14 {
		return 2
	} else if activeDays >= 1 && activeDays <= 7 {
		return 1
	} else {
		return 0
	}
}
