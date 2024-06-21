package api

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
)

type JobUpdates struct {
	bg background.System
}

func NewJobUpdates(bg background.System) *JobStatus {
	return &JobStatus{bg: bg}
}

type JobStatusNotifier struct {
	ID string
	WS *websocket.Conn
}

//nolint:exhaustivestruct // field defaults will work
var upgrader = websocket.Upgrader{
	// We accept all requests. CORS is enforcement is expected to be
	// handled by a proxy in front of our service.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *JobUpdates) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade the http connection to the WebSocket protocol.
	c, err := upgrader.Upgrade(w, r, nil)

	// Unexpected error
	if err != nil {
		appErr := backend.NewErrActionFailed("update to websocket", err.Error())

		handleAppErr(w, appErr)

		return
	}

	// Make sure we close the websocket when done.
	defer closeWebsocket(c)

	sub := &JobStatusNotifier{ID: gonanoid.Must(), WS: c}

	h.bg.Subscribe(sub)

	// Set up a defer function to unsubscribe
	defer h.bg.Unsubscribe(sub.ID)

	for {
		if _, _, err := c.NextReader(); err != nil {
			break
		}
	}
}

func (n *JobStatusNotifier) GetID() string {
	return n.ID
}

func (n *JobStatusNotifier) OnUpdate(s background.Status) {
	if err := n.WS.WriteJSON(s); err != nil {
		log.Error().Err(err).Msg("failed to write job status to websocket")
	}
}

func closeWebsocket(c io.Closer) {
	log.Info().Msg("Closing websocket")

	err := c.Close()
	if err != nil {
		log.Error().Err(err).Msg("could not close websocket")
	}
}
