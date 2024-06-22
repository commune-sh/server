package app

import (
	"context"
	"encoding/json"
	"net/http"
)

type SearchCategories struct {
	SearchCategories struct {
		RoomEvents struct {
			SearchTerm string `json:"search_term"`
			OrderBy    string `json:"order_by"`
		} `json:"room_events"`
	} `json:"search_categories"`
}

type SearchResults struct {
	SearchCategories struct {
		RoomEvents struct {
			SearchTerm string `json:"search_term"`
			OrderBy    string `json:"order_by"`
			Results    []struct {
				Rank   float64 `json:"rank"`
				Result struct {
					Type           string          `json:"type"`
					Sender         string          `json:"sender"`
					Content        json.RawMessage `json:"content"`
					StateKey       string          `json:"state_key"`
					OriginServerTS int64           `json:"origin_server_ts"`
					Unsigned       json.RawMessage `json:"unsigned"`
					EventID        string          `json:"event_id"`
					RoomID         string          `json:"room_id"`
				} `json:"result"`
				Context json.RawMessage `json:"context,omitempty"`
			} `json:"results"`
		} `json:"room_events"`
	} `json:"search_categories"`
}

func (c *App) SearchRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p, err := ReadRequestJSON(r, w, &SearchCategories{})

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"errcode": "M_BAD_JSON",
					"error":   "Content not JSON.",
				},
			})
			return
		}

		results, err := c.MatrixDB.Queries.SearchPublicRoomMessages(context.Background(), p.SearchCategories.RoomEvents.SearchTerm)
		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: results,
		})

	}
}
