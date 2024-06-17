package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	matrix_db "commune/db/matrix/gen"

	"github.com/go-chi/chi/v5"
	"github.com/tidwall/gjson"
)

type Event struct {
	Type           string          `json:"type"`
	Sender         string          `json:"sender"`
	Content        json.RawMessage `json:"content"`
	OriginServerTS int64           `json:"origin_server_ts"`
	Unsigned       json.RawMessage `json:"unsigned,omitempty"`
	EventID        string          `json:"event_id"`
	RoomID         string          `json:"room_id"`
}

func (c *App) RoomMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var limit int64 = 50

		l := r.URL.Query().Get("limit")
		lp, _ := strconv.ParseInt(l, 10, 64)
		if lp > 0 {
			limit = lp
		}

		dir := r.URL.Query().Get("dir")
		//from := r.URL.Query().Get("from")

		room_id := chi.URLParam(r, "room_id")

		rmp := matrix_db.GetRoomMessagesParams{
			RoomID: room_id,
			Limit:  &limit,
		}

		if dir == "b" {
			rmp.OrderBy = "DESC"
		} else if dir == "f" {
			rmp.OrderBy = "ASC"
		}

		messages, err := c.MatrixDB.Queries.GetRoomMessages(context.Background(), rmp)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"errorcode": "M_NOT_FOUND",
					"error":     err.Error(),
				},
			})
			return
		}

		processed, err := c.ProcessEvents(messages)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"errorcode": "M_UNKNOWN",
					"error":     err.Error(),
				},
			})
			return
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"messages": processed,
			},
		})

	}
}

func (c *App) ProcessEvents(events []matrix_db.GetRoomMessagesRow) (*[]Event, error) {

	processed := []Event{}

	for _, event := range events {
		e := Event{
			EventID: event.EventID,
		}

		content := gjson.Get(event.JSON, "content")
		if content.String() != "" {
			e.Content = json.RawMessage(content.Raw)
		}

		typ := gjson.Get(event.JSON, "type")
		if typ.String() != "" {
			e.Type = typ.String()
		}

		evid := gjson.Get(event.JSON, "event_id")
		if evid.String() != "" {
			e.EventID = evid.String()
		}

		rid := gjson.Get(event.JSON, "room_id")
		if rid.String() != "" {
			e.RoomID = rid.String()
		}

		sender := gjson.Get(event.JSON, "sender")
		if sender.String() != "" {
			e.Sender = sender.String()
		}

		origin_server_ts := gjson.Get(event.JSON, "origin_server_ts")
		if origin_server_ts.String() != "" {
			e.OriginServerTS = origin_server_ts.Int()
		}

		unsigned := gjson.Get(event.JSON, "unsigned")
		if unsigned.String() != "" {
			e.Unsigned = json.RawMessage(unsigned.Raw)
		}

		processed = append(processed, e)
	}

	return &processed, nil
}
