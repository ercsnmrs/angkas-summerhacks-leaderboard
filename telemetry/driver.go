package telemetry

import (
	"context"
	"fmt"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type DriverServiceInstrumentation struct {
	serviceName string
	svc         *driver.Service
}

func NewDriverServiceInstrumentation(
	serviceName string,
	svc driver.Service,
) *DriverServiceInstrumentation {
	return &DriverServiceInstrumentation{
		serviceName: fmt.Sprintf("incentive-%s-service", serviceName),
		svc:         &svc,
	}
}

func (t *DriverServiceInstrumentation) GetDriver(ctx context.Context, driverID string) (rating driver.Driver, err error) {
	ctx, span := otel.Tracer(t.serviceName).Start(ctx, "get-driver-rating")
	defer span.End()

	r, err := t.svc.GetDriver(ctx, driverID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return rating, err
	}

	return r, nil
}

func (t *DriverServiceInstrumentation) WriteToCSV(ctx context.Context, userID string, rating driver.RFM, driver string) error {
	ctx, span := otel.Tracer(t.serviceName).Start(ctx, "write-to-csv")
	defer span.End()

	err := t.svc.WriteToCSV(ctx, userID, rating, driver)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// func (t *DriverServiceInstrumentation) UpdateUserRating(ctx context.Context, trip ) error {
// 	ctx, span := otel.Tracer(t.serviceName).Start(ctx, "write-to-csv")
// 	defer span.End()

// 	err := t.svc.UpdateUserRating(ctx, userID, )
// 	if err != nil {
// 		span.RecordError(err)
// 		span.SetStatus(codes.Error, err.Error())
// 		return err
// 	}

// 	return nil
// }
