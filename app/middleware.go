package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// make sure this room exists
func (c *App) EnsureRoomExists(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		// missing room param
		if room_id == "" {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"error": "room ID is required",
				},
			})
			return
		}

		is_room_id := IsValidRoomID(room_id)

		if !is_room_id {
			room_id = fmt.Sprintf("!%s:%s", room_id, c.Config.Matrix.ServerName)
			rctx := chi.RouteContext(r.Context())
			rctx.URLParams.Add("room_id", room_id)
		}

		log.Println("room_id:", room_id)

		// does room exist?
		exists, err := c.MatrixDB.Queries.DoesRoomExist(context.Background(), room_id)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		if !exists {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"errcode": "ROOM_NOT_FOUND",
					"error":   "room does not exist",
				},
			})
			return
		}

		h.ServeHTTP(w, r)
	})
}

// validate room is public
func (c *App) ValidateRoomIsPublic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

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
				Code: http.StatusForbidden,
				JSON: map[string]any{
					"errcode": "ROOM_NOT_PUBLIC",
					"error":   "room is not public",
				},
			})
			return
		}

		h.ServeHTTP(w, r)
	})
}
