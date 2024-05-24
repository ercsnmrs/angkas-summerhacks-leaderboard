package trip

import (
	"fmt"
	"math/rand"
	"time"
)

type TripCustomerInfo struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type TripPlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

type Address struct {
	Main      string `json:"main"`
	Secondary string `json:"secondary"`
}

type TripLocation struct {
	Latitude      string       `json:"latitude"`
	Longitude     string       `json:"longitude"`
	PlusCode      TripPlusCode `json:"plus_code"`
	Address       Address      `json:"address"`
	AngkasPlaceID string       `json:"angkas_place_id"`
}

type DriverInfo struct {
	Rating               int       `json:"rating"`
	PickupDistanceMeters int       `json:"pickup_distance_meters"`
	PickupETA            time.Time `json:"pickup_eta"`
}

type PriceInfo struct {
	Total                int     `json:"total"`
	Final                int     `json:"final"`
	Surcharge            int     `json:"surcharge"`
	BasePrice            int     `json:"base_price"`
	TimeBasedPrice       int     `json:"time_based_price"`
	TimeBasedPriceRate   int     `json:"time_based_price_rate"`
	Discount             int     `json:"discount"`
	Currency             string  `json:"currency"`
	CommissionPercentage int     `json:"commission_percentage"`
	DriverEarnings       float64 `json:"driver_earnings"`
	ServiceFee           int     `json:"service_fee"`
	ReferenceID          string  `json:"reference_id"`
	Distance             int     `json:"distance"`
	Duration             int     `json:"duration"`
	Units                int     `json:"units"`
}

type PaymentInfo struct {
	ID                string `json:"id"`
	PaymentMethodID   string `json:"payment_method_id"`
	PaymentMethodType string `json:"payment_method_type"`
	Amount            int    `json:"amount"`
}

type MetadataInfo struct {
	RoutesMethod              string    `json:"routes_method"`
	TripDistance              int       `json:"trip_distance"`
	PickupDistance            int       `json:"pickup_distance"`
	DriverDistanceToDropoff   int       `json:"driver_distance_to_dropoff"`
	PickupEstimatedTimestamp  time.Time `json:"pickup_estimated_timestamp"`
	DropoffEstimatedTimestamp time.Time `json:"dropoff_estimated_timestamp"`
	BikerAcceptedTimestamp    time.Time `json:"biker_accepted_timestamp"`
	PickupTimestamp           time.Time `json:"pickup_timestamp"`
	DropoffTimestamp          time.Time `json:"dropoff_timestamp"`
	CompletedAt               time.Time `json:"completed_at"`
}

type Event struct {
	TripRequestID  string           `json:"trip_request_id"`
	OdrdTripID     string           `json:"odrd_trip_id"`
	CustomerID     string           `json:"customer_id"`
	CustomerInfo   TripCustomerInfo `json:"customer_info"`
	DriverID       string           `json:"driver_id"`
	ServiceID      string           `json:"service_id"`
	Pickup         TripLocation     `json:"pickup"`
	Dropoff        TripLocation     `json:"dropoff"`
	Driver         DriverInfo       `json:"driver"`
	Price          PriceInfo        `json:"price"`
	Status         string           `json:"status"`
	Notes          string           `json:"notes"`
	ServiceZone    string           `json:"service_zone"`
	Payment        PaymentInfo      `json:"payment"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Feedback       interface{}      `json:"feedback"`
	Metadata       MetadataInfo     `json:"metadata"`
	IdempotencyKey string           `json:"idempotency_key"`
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}

func randomPhoneNumber() string {
	return fmt.Sprintf("+%d%d%d%d%d%d%d%d%d%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}

func randomEmail() string {
	return fmt.Sprintf("%s@example.com", randomString(10))
}

func randomLatitude() string {
	return fmt.Sprintf("%.6f", (rand.Float64()*180)-90)
}

func randomLongitude() string {
	return fmt.Sprintf("%.6f", (rand.Float64()*360)-180)
}

func randomAddress() Address {
	return Address{
		Main:      randomString(10) + " Street",
		Secondary: randomString(5) + " Building",
	}
}

func randomTripLocation() TripLocation {
	return TripLocation{
		Latitude:  randomLatitude(),
		Longitude: randomLongitude(),
		PlusCode: TripPlusCode{
			CompoundCode: randomString(8),
			GlobalCode:   randomString(10),
		},
		Address:       randomAddress(),
		AngkasPlaceID: randomString(8),
	}
}

func randomTimeInRange(start, end time.Time) time.Time {
	delta := end.Sub(start)
	sec := rand.Int63n(int64(delta.Seconds()))
	return start.Add(time.Duration(sec) * time.Second)
}

func GenerateFakeTripEvent() Event {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	createdAt := randomTimeInRange(now.Add(-30*24*time.Hour), now)
	updatedAt := randomTimeInRange(createdAt, now)
	event := Event{
		TripRequestID: randomString(10),
		OdrdTripID:    randomString(10),
		CustomerID:    randomString(10),
		CustomerInfo: TripCustomerInfo{
			Name:        randomString(5) + " " + randomString(7),
			PhoneNumber: randomPhoneNumber(),
			Email:       randomEmail(),
		},
		DriverID:       randomString(10),
		ServiceID:      randomString(10),
		Pickup:         randomTripLocation(),
		Dropoff:        randomTripLocation(),
		Driver:         DriverInfo{Rating: rand.Intn(5) + 1, PickupDistanceMeters: rand.Intn(10000), PickupETA: randomTimeInRange(now.Add(-1*time.Hour), now)},
		Price:          PriceInfo{Total: rand.Intn(5000), Final: rand.Intn(5000), Surcharge: rand.Intn(500), BasePrice: rand.Intn(3000), TimeBasedPrice: rand.Intn(2000), TimeBasedPriceRate: rand.Intn(50), Discount: rand.Intn(500), Currency: "USD", CommissionPercentage: rand.Intn(100), DriverEarnings: rand.Float64(), ServiceFee: rand.Intn(200), ReferenceID: randomString(8), Distance: rand.Intn(100), Duration: rand.Intn(100), Units: rand.Intn(2)},
		Status:         "complete",
		Notes:          randomString(20),
		ServiceZone:    randomString(10),
		Payment:        PaymentInfo{ID: randomString(8), PaymentMethodID: randomString(8), PaymentMethodType: "card", Amount: rand.Intn(5000)},
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Feedback:       nil,
		Metadata:       MetadataInfo{RoutesMethod: randomString(10), TripDistance: rand.Intn(500), PickupDistance: rand.Intn(100), DriverDistanceToDropoff: rand.Intn(200), PickupEstimatedTimestamp: randomTimeInRange(createdAt, updatedAt), DropoffEstimatedTimestamp: randomTimeInRange(createdAt, updatedAt), BikerAcceptedTimestamp: randomTimeInRange(createdAt, updatedAt), PickupTimestamp: randomTimeInRange(createdAt, updatedAt), DropoffTimestamp: randomTimeInRange(createdAt, updatedAt), CompletedAt: randomTimeInRange(createdAt, updatedAt)},
		IdempotencyKey: randomString(10),
	}

	return event
}
