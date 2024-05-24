package trip

// func TestGenerateFakeEvents(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	now := time.Now()

// 	var events []Event
// 	for i := 0; i < 1000; i++ {
// 		createdAt := randomTimeInRange(now.Add(-30*24*time.Hour), now)
// 		updatedAt := randomTimeInRange(createdAt, now)
// 		event := Event{
// 			TripRequestID: randomString(10),
// 			OdrdTripID:    randomString(10),
// 			CustomerID:    randomString(10),
// 			CustomerInfo: TripCustomerInfo{
// 				Name:        randomString(5) + " " + randomString(7),
// 				PhoneNumber: randomPhoneNumber(),
// 				Email:       randomEmail(),
// 			},
// 			DriverID:       randomString(10),
// 			ServiceID:      randomString(10),
// 			Pickup:         randomTripLocation(),
// 			Dropoff:        randomTripLocation(),
// 			Driver:         DriverInfo{Rating: rand.Intn(5) + 1, PickupDistanceMeters: rand.Intn(10000), PickupETA: randomTimeInRange(now.Add(-1*time.Hour), now)},
// 			Price:          PriceInfo{Total: rand.Intn(5000), Final: rand.Intn(5000), Surcharge: rand.Intn(500), BasePrice: rand.Intn(3000), TimeBasedPrice: rand.Intn(2000), TimeBasedPriceRate: rand.Intn(50), Discount: rand.Intn(500), Currency: "USD", CommissionPercentage: rand.Intn(100), DriverEarnings: rand.Intn(3000), ServiceFee: rand.Intn(200), ReferenceID: randomString(8), Distance: rand.Intn(100), Duration: rand.Intn(100), Units: rand.Intn(2)},
// 			Status:         "complete",
// 			Notes:          randomString(20),
// 			ServiceZone:    randomString(10),
// 			Payment:        PaymentInfo{ID: randomString(8), PaymentMethodID: randomString(8), PaymentMethodType: "card", Amount: rand.Intn(5000)},
// 			CreatedAt:      createdAt,
// 			UpdatedAt:      updatedAt,
// 			Feedback:       nil,
// 			Metadata:       MetadataInfo{RoutesMethod: randomString(10), TripDistance: rand.Intn(500), PickupDistance: rand.Intn(100), DriverDistanceToDropoff: rand.Intn(200), PickupEstimatedTimestamp: randomTimeInRange(createdAt, updatedAt), DropoffEstimatedTimestamp: randomTimeInRange(createdAt, updatedAt), BikerAcceptedTimestamp: randomTimeInRange(createdAt, updatedAt), PickupTimestamp: randomTimeInRange(createdAt, updatedAt), DropoffTimestamp: randomTimeInRange(createdAt, updatedAt), CompletedAt: randomTimeInRange(createdAt, updatedAt)},
// 			IdempotencyKey: randomString(10),
// 		}
// 		events = append(events, event)
// 	}

// 	// Convert to JSON and print the events
// 	// eventsJSON, err := json.MarshalIndent(events, "", "  ")
// 	// if err != nil {
// 	// 	fmt.Println("Error marshalling events to JSON:", err)
// 	// 	return
// 	// }

// 	// fmt.Println(string(eventsJSON))
// }
