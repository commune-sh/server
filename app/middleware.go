package app

import (
	"context"
	"fmt"
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

			// if this is a room alias, let's see if it belongs to a valid room

			alias := fmt.Sprintf("#%s:%s", room_id, c.Config.Matrix.ServerName)

			id, err := c.MatrixDB.Queries.GetRoomIDFromAlias(context.Background(), alias)

			if err == nil && id != "" {
				room_id = id
			} else {
				room_id = fmt.Sprintf("!%s:%s", room_id, c.Config.Matrix.ServerName)
			}

			rctx := chi.RouteContext(r.Context())
			rctx.URLParams.Add("room_id", room_id)
		}

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
		//is_public, err := c.MatrixDB.Queries.IsRoomPublic(context.Background(), room_id)
		is_public, err := c.MatrixDB.Queries.IsRoomPubliclyAccessible(context.Background(), room_id)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		c.Log.Debug().
			Bool("Public", is_public).
			Msg("Is room public?")

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

// makes sure this route is autehnticated
func (c *App) RequireAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access_token, err := ExtractAccessToken(r)

		if err != nil || access_token == nil {
			RespondWithJSON(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"errcode": "BAD_ACCESS_TOKEN",
					"error":   "access token invalid",
				},
			})
			return
		}

		v, err := c.MatrixDB.Queries.IsAccessTokenValid(context.Background(), *access_token)

		if err != nil || !v.Valid || v.UserID == "" {
			RespondWithJSON(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"errcode":     "M_UNKNOWN_TOKEN",
					"error":       "Invalid access token passed.",
					"soft_logout": false,
				},
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", v.UserID)
		ctx = context.WithValue(ctx, "access_token", *access_token)

		h.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (c *App) AuthenticatedUser(r *http.Request) *string {
	user_id, ok := r.Context().Value("user_id").(string)

	if !ok {
		return nil
	}

	return &user_id
}

func (c *App) AuthenticatedAccessToken(r *http.Request) *string {
	access_token, ok := r.Context().Value("access_token").(string)

	if !ok {
		return nil
	}

	return &access_token

}

// makes sure this route is autehnticated
func (c *App) RequireAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user_id := c.AuthenticatedUser(r)

		is_admin, err := c.MatrixDB.Queries.IsUserAdmin(context.Background(), user_id)

		if err != nil || !is_admin {
			RespondWithJSON(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"errcode": "M_FORBIDDEN",
					"error":   "You are not a server admin",
				},
			})
			return
		}

		h.ServeHTTP(w, r)

	})
}
