package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/leaderboard"
)

type leaderboardService interface {
	GetLeaderboard(ctx context.Context, scope string) (drivers leaderboard.Leaderboard, err error)
}

func GetLeaderboard(svc leaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		scope := chi.URLParam(r, "scope")

		c, err := svc.GetLeaderboard(r.Context(), scope)
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
