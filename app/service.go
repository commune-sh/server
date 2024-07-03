package app

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"maunium.net/go/mautrix/event"
)

func (c *App) Transactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var events struct {
			Events []event.Event `json:"events"`
		}
		if err := json.Unmarshal(body, &events); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, event := range events.Events {
			if event.RoomID != "" {
				c.Log.Info().Msgf("Event: %v", event)
			}

			join, err := c.Matrix.JoinRoom(context.Background(), event.RoomID.String(), "", nil)

			if err != nil {
				c.Log.Error().Msgf("Error joining room: %v", err)
			}

			c.Log.Info().Msgf("Join response: %v", join)

		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
		})

	}
}
