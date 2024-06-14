package app

import (
	"net/http"
)

func (c *App) Capabilities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"capabilities": c.Config.Capabilities,
			},
		})

	}
}
