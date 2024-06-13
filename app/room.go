package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tidwall/gjson"
)

type RoomHierarchyItem struct {
	RoomID           string        `json:"room_id"`
	RoomType         string        `json:"room_type"`
	CanonicalAlias   string        `json:"canonical_alias,omitempty"`
	JoinRule         string        `json:"join_rule"`
	WorldReadable    bool          `json:"world_readable"`
	Name             string        `json:"name"`
	Topic            string        `json:"topic,omitempty"`
	NumJoinedMembers int64         `json:"num_joined_members"`
	GuestCanJoin     bool          `json:"guest_can_join"`
	AvatarURL        string        `json:"avatar_url,omitempty"`
	ChildrenState    []interface{} `json:"children_state,omitempty"`
}

func (c *App) RoomHierarchy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		h, err := c.MatrixDB.Queries.GetRoomHierarchy(context.Background(), room_id)
		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		rooms := []RoomHierarchyItem{}

		if len(h) > 0 {
			for _, id := range h {

				room := RoomHierarchyItem{
					RoomID: id,
				}

				events, err := c.MatrixDB.Queries.GetCurrentStateEvents(context.Background(), id)
				if err != nil {
					continue
				}

				for _, x := range events {

					if x.CurrentStateEvent == "m.room.create" {
						item := gjson.Get(x.Content.String, "content.type")
						if item.String() != "" {
							room.RoomType = item.String()
						}
					}

					if x.CurrentStateEvent == "m.room.name" {
						item := gjson.Get(x.Content.String, "content.name")
						if item.String() != "" {
							room.Name = item.String()
						}
					}

					if x.CurrentStateEvent == "m.room.canonical_alias" {
						item := gjson.Get(x.Content.String, "content.alias")
						if item.String() != "" {
							room.CanonicalAlias = item.String()
						}
					}

					if x.CurrentStateEvent == "m.room.join_rules" {
						item := gjson.Get(x.Content.String, "content.join_rule")
						if item.String() != "" {
							room.JoinRule = item.String()
						}
					}

					if x.CurrentStateEvent == "m.room.history_visibility" {
						item := gjson.Get(x.Content.String, "content.history_visibility")
						if item.String() == "world_readable" {
							room.WorldReadable = true
						}
					}

					if x.CurrentStateEvent == "m.room.guest_access" {
						item := gjson.Get(x.Content.String, "content.guest_access")
						if item.String() == "can_join" {
							room.GuestCanJoin = true
						}
					}

					if x.CurrentStateEvent == "m.room.avatar" {
						item := gjson.Get(x.Content.String, "content.url")
						if item.String() != "" {
							room.AvatarURL = item.String()
						}
					}

					if x.CurrentStateEvent == "m.space.child" {

						type ChildState struct {
							Type           string          `json:"type"`
							StateKey       string          `json:"state_key"`
							Content        json.RawMessage `json:"content"`
							Sender         string          `json:"sender"`
							OriginServerTS int64           `json:"origin_server_ts"`
						}

						cs := ChildState{}

						content := gjson.Get(x.Content.String, "content")
						if content.String() != "" {
							cs.Content = json.RawMessage(content.Raw)
						}

						typ := gjson.Get(x.Content.String, "type")
						if typ.String() != "" {
							cs.Type = typ.String()
						}

						state_key := gjson.Get(x.Content.String, "state_key")
						if state_key.String() != "" {
							cs.StateKey = state_key.String()
						}

						sender := gjson.Get(x.Content.String, "sender")
						if sender.String() != "" {
							cs.Sender = sender.String()
						}

						origin_server_ts := gjson.Get(x.Content.String, "origin_server_ts")
						if origin_server_ts.String() != "" {
							cs.OriginServerTS = origin_server_ts.Int()
						}

						room.ChildrenState = append(room.ChildrenState, cs)
					}
				}

				rooms = append(rooms, room)
			}
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"rooms": rooms,
			},
		})

	}
}
