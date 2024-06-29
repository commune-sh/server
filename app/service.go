package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

		authBearer := r.Header.Get("Authorization")

		log.Println("Authorization: ", authBearer)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var events struct {
			Events []MatrixEvent `json:"events"`
		}
		if err := json.Unmarshal(body, &events); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, event := range events.Events {
			c.Log.Info().Msgf("Event: %v", event)
		}

		w.WriteHeader(http.StatusOK)

	}
}
