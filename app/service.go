package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"maunium.net/go/mautrix/event"
)

type MatrixEvent struct {
	StateKey    *string                `json:"state_key,omitempty"`
	Sender      string                 `json:"sender"`
	Type        string                 `json:"type"`
	Timestamp   int64                  `json:"origin_server_ts"`
	ID          string                 `json:"event_id"`
	RoomID      string                 `json:"room_id"`
	Redacts     string                 `json:"redacts,omitempty"`
	Unsigned    map[string]interface{} `json:"unsigned"`
	Content     map[string]interface{} `json:"content"`
	PrevContent map[string]interface{} `json:"prev_content,omitempty"`
}

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
			c.Log.Info().Msgf("Event: %v", event)
		}

		RespondWithJSON(w, &JSONResponse{
			Code: http.StatusOK,
		})

	}
}
