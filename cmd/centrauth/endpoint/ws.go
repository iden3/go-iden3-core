package endpoint

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type webSocketConn struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var (
	wsMutex    sync.Mutex
	wsClients  map[string]webSocketConn
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	wsClients = make(map[string]webSocketConn)
}

func handleWs(c *gin.Context) {
	w := c.Writer
	r := c.Request
	clientID := c.Param("id")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(w, r, nil)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	newClient := &webSocketConn{conn: conn, send: make(chan []byte, 256)}

	wsMutex.Lock()
	defer wsMutex.Unlock()
	wsClients[clientID] = *newClient
}

func wsWriteAndClose(clientID string, content []byte) {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	wsClients[clientID].conn.WriteMessage(websocket.TextMessage, content)
	delete(wsClients, clientID)
}
