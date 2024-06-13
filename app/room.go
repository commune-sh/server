package app

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/go-chi/chi/v5"
)

func (c *App) RoomHierarchy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		room_id := chi.URLParam(r, "room_id")

		h, err := c.MatrixDB.Queries.GetRoomHierarchy(context.Background(), room_id)

		type room struct {
			RoomID string `json:"room_id"`
			State  any    `json:"state"`
		}

		rooms := []room{}

		if len(h) > 0 {
			for _, id := range h {

				events, err := c.MatrixDB.Queries.GetCurrentStateEvents(context.Background(), id)
				if err != nil {
					continue
				}

				for _, x := range events {

					doc, err := jsonquery.Parse(strings.NewReader(x.Content.String))
					if err != nil {
						log.Println(err)
					}
					age := jsonquery.FindOne(doc, "content/name")
					if age != nil {
						log.Println(age.Value())
					}

				}

				rooms = append(rooms, room{
					RoomID: id,
					State:  events,
				})
			}
		}

		if err != nil {
			RespondWithJSON(w, &JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
			JSON: map[string]any{
				"rooms": rooms,
			},
		})

	}
}

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
