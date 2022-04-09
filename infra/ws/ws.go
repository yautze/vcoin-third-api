package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader -
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebSocket -
type WebSocket struct {
	mux               sync.RWMutex
	Conn              *websocket.Conn
	Out               chan []byte
	OutError          chan *Event
	In                chan *Event
	Closer            chan bool
	Events            map[string]EventHandler
	EventsWithChannel map[string]chan *Event
}

// NewWebSocket -
func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	subprotocols := r.Header["Sec-Websocket-Protocol"]
	upgrader.Subprotocols = subprotocols

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("An error occured while upgrading the connection: %v", err)
		return nil, err
	}

	ws := &WebSocket{
		Conn:              conn,
		Out:               make(chan []byte),
		OutError:          make(chan *Event),
		In:                make(chan *Event),
		Closer:            make(chan bool),
		Events:            make(map[string]EventHandler),
		EventsWithChannel: make(map[string]chan *Event),
	}

	go ws.Reader()
	go ws.Writer()

	return ws, nil
}

// Reader -
func (ws *WebSocket) Reader() {
	// close websocket connection
	defer ws.ConnectClose()

	for {
		_, message, err := ws.Conn.ReadMessage()

		// reader err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("webSocket error")
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("Client reload disconnected")
				return
			}

			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				log.Println("Client actively disconnected")
				return
			}

			log.Println("other error")

			break
		}
		event, err := NewEvent(message)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
		}

		if action, ok := ws.Events[event.Name]; ok {
			action(event)
		}

		if action, ok := ws.EventsWithChannel[event.Name]; ok {
			action <- event
		}
	}
}

// Writer -
func (ws *WebSocket) Writer() {
	defer ws.ConnectClose()

	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				ws.Conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
			}

			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			ws.mux.Lock()
			w.Write(message)
			w.Close()
			ws.mux.Unlock()

		case event := <-ws.OutError:
			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			r := make(map[string]interface{})
			r["time"] = time.Now().Unix()
			r["code"] = 21000
			r["message"] = event.Data.(string)

			res, _ := json.Marshal(r)

			event.Data = string(res)

			w.Write(event.Raw())
			w.Close()

			return
		}
	}
}

// On -
func (ws *WebSocket) On(eventName string, action EventHandler) *WebSocket {
	ws.Events[eventName] = action
	return ws
}

// OnChannel -
func (ws *WebSocket) OnChannel(eventName string, ch chan *Event) chan *Event {
	ws.EventsWithChannel[eventName] = ch
	return ch
}

// ConnectClose -
func (ws *WebSocket) ConnectClose() {
	// send close to channel
	ws.Closer <- true

	// close connection
	ws.Conn.Close()
}
