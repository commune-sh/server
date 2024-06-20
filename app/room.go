package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	NumJoinedMembers int32         `json:"num_joined_members"`
	GuestCanJoin     bool          `json:"guest_can_join"`
	AvatarURL        string        `json:"avatar_url,omitempty"`
	ChildrenState    []interface{} `json:"children_state,omitempty"`
}

func (c *App) RoomHierarchy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		//suggested_only := r.URL.Query().Get("suggested_only")
		//limit := r.URL.Query().Get("limit")

		// get room hierarchy
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
			for _, state := range h {

				room := RoomHierarchyItem{
					RoomID: state.RoomID,
				}

				if state.RoomType != nil {
					room.RoomType = *state.RoomType
				}

				if state.Name != nil {
					room.Name = *state.Name
				}

				if state.Topic != nil {
					room.Topic = *state.Topic
				}

				if state.CanonicalAlias != nil {
					room.CanonicalAlias = *state.CanonicalAlias
				}

				if state.JoinRules != nil {
					room.JoinRule = *state.JoinRules
				}

				if state.HistoryVisibility != nil && *state.HistoryVisibility == "world_readable" {
					room.WorldReadable = true
				}

				if state.GuestAccess != nil && *state.GuestAccess == "can_join" {
					room.GuestCanJoin = true
				}

				if state.Avatar != nil {
					room.AvatarURL = *state.Avatar
				}

				if state.JoinedMembers > 0 {
					room.NumJoinedMembers = state.JoinedMembers
				}

				if state.IsParent {

					// get current state events
					events, err := c.MatrixDB.Queries.GetSpaceChildStateEvents(context.Background(), state.RoomID)
					if err != nil {
						continue
					}

					for _, x := range events {

						if x.CurrentStateEvent == "m.space.child" {

							type ChildState struct {
								Type           string          `json:"type"`
								StateKey       string          `json:"state_key"`
								Content        json.RawMessage `json:"content"`
								Sender         string          `json:"sender"`
								OriginServerTS int64           `json:"origin_server_ts"`
							}

							cs := ChildState{}

							content := gjson.Get(x.EventJson, "content")
							if content.String() != "" {
								cs.Content = json.RawMessage(content.Raw)
							}

							typ := gjson.Get(x.EventJson, "type")
							if typ.String() != "" {
								cs.Type = typ.String()
							}

							state_key := gjson.Get(x.EventJson, "state_key")
							if state_key.String() != "" {
								cs.StateKey = state_key.String()
							}

							sender := gjson.Get(x.EventJson, "sender")
							if sender.String() != "" {
								cs.Sender = sender.String()
							}

							origin_server_ts := gjson.Get(x.EventJson, "origin_server_ts")
							if origin_server_ts.String() != "" {
								cs.OriginServerTS = origin_server_ts.Int()
							}

							room.ChildrenState = append(room.ChildrenState, cs)
						}
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

func (c *App) PublicRooms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var limit int64 = 10

		l := r.URL.Query().Get("limit")
		lp, _ := strconv.ParseInt(l, 10, 64)
		if lp > 0 {
			limit = lp
		}

		spaces, err := c.MatrixDB.Queries.GetPublicSpaces(context.Background(), &limit)
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

		if len(spaces) > 0 {
			for _, state := range spaces {

				room := RoomHierarchyItem{
					RoomID: state.RoomID,
				}

				if state.JoinedMembers > 0 {
					room.NumJoinedMembers = state.JoinedMembers
				}

				if state.RoomType != nil {
					room.RoomType = *state.RoomType
				}

				if state.Name != nil {
					room.Name = *state.Name
				}

				if state.Topic != nil {
					room.Topic = *state.Topic
				}

				if state.CanonicalAlias != nil {
					room.CanonicalAlias = *state.CanonicalAlias
				}

				if state.JoinRules != nil {
					room.JoinRule = *state.JoinRules
				}

				if state.HistoryVisibility != nil && *state.HistoryVisibility == "world_readable" {
					room.WorldReadable = true
				}

				if state.GuestAccess != nil && *state.GuestAccess == "can_join" {
					room.GuestCanJoin = true
				}

				if state.Avatar != nil {
					room.AvatarURL = *state.Avatar
				}

				rooms = append(rooms, room)
			}
		}

		total, _ := c.MatrixDB.Queries.GetPublicSpacesCount(context.Background())

		resp := map[string]any{
			"chunk":                     rooms,
			"total_room_count_estimate": total,
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: resp,
		})

	}
}

type StateEvent struct {
	Type           string          `json:"type"`
	Sender         string          `json:"sender"`
	Content        json.RawMessage `json:"content"`
	StateKey       string          `json:"state_key"`
	OriginServerTS int64           `json:"origin_server_ts"`
	Unsigned       json.RawMessage `json:"unsigned,omitempty"`
	EventID        string          `json:"event_id"`
	RoomID         string          `json:"room_id"`
}

func (c *App) RoomStateEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		events, err := c.MatrixDB.Queries.GetCurrentStateEvents(context.Background(), room_id)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		cse := []StateEvent{}

		if len(events) > 0 {

			for _, x := range events {

				cs := StateEvent{
					RoomID: room_id,
				}

				if x.EventID != "" {
					cs.EventID = x.EventID
				}

				content := gjson.Get(x.EventJson, "content")
				if content.String() != "" {
					cs.Content = json.RawMessage(content.Raw)
				}

				typ := gjson.Get(x.EventJson, "type")
				if typ.String() != "" {
					cs.Type = typ.String()
				}

				state_key := gjson.Get(x.EventJson, "state_key")
				if state_key.String() != "" {
					cs.StateKey = state_key.String()
				}

				sender := gjson.Get(x.EventJson, "sender")
				if sender.String() != "" {
					cs.Sender = sender.String()
				}

				origin_server_ts := gjson.Get(x.EventJson, "origin_server_ts")
				if origin_server_ts.String() != "" {
					cs.OriginServerTS = origin_server_ts.Int()
				}

				unsigned := gjson.Get(x.EventJson, "unsigned")
				if unsigned.String() != "" {
					cs.Unsigned = json.RawMessage(unsigned.Raw)
				}

				cse = append(cse, cs)
			}
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"current_state_events": cse,
			},
		})

	}
}

func (c *App) RoomState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		state, err := c.MatrixDB.Queries.GetRoomState(context.Background(), room_id)

		if err != nil {
			RespondWithError(w, &JSONResponse{
				Code: http.StatusOK,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: json.RawMessage(state),
		})

	}
}

func (c *App) IsRoomPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		is_public, err := c.MatrixDB.Queries.IsRoomPubliclyAccessible(context.Background(), room_id)

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

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"public": is_public,
			},
		})

	}
}
