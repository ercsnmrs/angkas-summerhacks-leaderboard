package faker

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
)

func GenerateFakeDrivers(numDrivers int) []driver.Driver {
	var drivers []driver.Driver
	for i := 0; i < numDrivers; i++ {
		dr := driver.Driver{
			DriverID:                     uuid.New().String(),
			LastCompletedTripDate:        generateRandomDate(),
			NetIncome:                    generateRandomFloat(1000, 50000),
			NumberOfCompletedTrips:       generateRandomInt(1, 500),
			UniqueDateWithCompletedTrips: generateRandomInt(1, 30),
			ServiceZone:                  generateRandomServiceZone(),
			Rating: driver.Rating{
				RFM: driver.RFM{
					Recency:   generateRandomFloat(1, 5),
					Frequency: generateRandomFloat(1, 5),
					Monetary:  generateRandomFloat(1, 5),
				},
			},
		}

		dr.Rating.Average = dr.CalculateAverage()

		drivers = append(drivers, dr)
	}
	return drivers
}

func generateRandomDate() time.Time {
	// Generate a random number of days between 0 and 30
	randomDays := rand.Intn(30)
	// Subtract the random number of days from the current time
	lastCompletedTripDate := time.Now().AddDate(0, 0, -randomDays)
	return lastCompletedTripDate
}

func generateRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func generateRandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func generateRandomServiceZone() string {
	serviceZones := []string{"MNL", "CEB", "CDO"}
	randomIndex := rand.Intn(len(serviceZones))
	return serviceZones[randomIndex]
}
