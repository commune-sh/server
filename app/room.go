package app

import (
	"net/http"
)

func (c *App) Test() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"message": "test",
			},
		})

	}
}
