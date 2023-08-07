package websocket

import (
	"bytes"
	"log"
	"net/http"
	roomService "sirkelin/backend/app/room/service"
	"sirkelin/backend/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan models.Message
	roomService *roomService.RoomService
	roomID      string
	userID      string
}

type IClient interface {
	readPump()
	writePump()
}

func NewClient(hub *Hub, conn *websocket.Conn, send chan models.Message, roomService *roomService.RoomService, roomID, userID string) *Client {
	return &Client{
		hub:         hub,
		conn:        conn,
		send:        send,
		roomService: roomService,
		roomID:      roomID,
		userID:      userID,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		stringifyMessage := string(message[:])
		messageID, _ := c.roomService.PushMessage(c.userID, c.roomID, models.SendMessageParams{Body: stringifyMessage})
		c.hub.broadcast <- models.Message{
			ID:        messageID,
			Body:      stringifyMessage,
			UserID:    c.userID,
			RoomID:    c.roomID,
			CreatedAt: time.Now(),
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteJSON(message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, roomService *roomService.RoomService, roomID, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	client := NewClient(hub, conn, make(chan models.Message), roomService, roomID, userID)
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
