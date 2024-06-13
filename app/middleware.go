package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// validate if room exists and is public
func (c *App) ValidateRoom(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		// missing room param
		if room_id == "" {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": "room_id is required",
				},
			})
			return
		}

		// does room exist?
		exists, err := c.MatrixDB.Queries.DoesRoomExist(context.Background(), room_id)

		if err != nil || !exists {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": "room does not exist",
				},
			})
			return
		}

		// is room public?
		is_public, err := c.MatrixDB.Queries.IsRoomPublic(context.Background(), room_id)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		// not public?
		if !is_public {
			RespondWithJSON(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": "room is not public",
				},
			})
			return
		}

		h.ServeHTTP(w, r)
	})
}
