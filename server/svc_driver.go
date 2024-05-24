package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
)

type driverService interface {
	GetDriver(ctx context.Context, id string) (driver driver.Driver, err error)
}

func GetDriverRating(ds driverService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := chi.URLParam(r, "id")
		c, err := ds.GetDriver(r.Context(), id)
		if err != nil {
			encodeJSONError(w, err, http.StatusBadRequest)
			return
		}

		encodeJSONResp(w, c, http.StatusOK)
	}
}

// func ListTier(s tierService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		encodeJSONResp(w, struct {
// 			Msg string `json:"message"`
// 		}{"no fighters yet implemented"}, http.StatusOK)
// 	}
// }
