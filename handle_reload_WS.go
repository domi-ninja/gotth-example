package app

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func (webapp *WebApp) HandleReload_WS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	conn.WriteMessage(websocket.TextMessage, []byte(webapp.version))

	for {
		// keep the connection open while we are running
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}
