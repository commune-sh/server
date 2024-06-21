package app

import (
	"context"
	"net/http"
)

type Space struct {
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

func (c *App) GetUserSpaces() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user_id := c.AuthenticatedUser(r)

		spaces, err := c.MatrixDB.Queries.GetUserSpaces(context.Background(), user_id)

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

		rooms := []RoomItem{}

		if len(spaces) > 0 {
			for _, state := range spaces {

				room := RoomItem{
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

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: rooms,
		})

	}
}
