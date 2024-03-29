package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"service.chat/internal/logic"
	"service.chat/internal/svc"
	"service.chat/internal/types"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 3) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// send buffer size
	bufSize = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	svc *svc.ServiceContext
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	UserId string
	JwtId  string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(50240)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		le, message, err := c.conn.ReadMessage()

		fmt.Println("+++", le, len(message), err)
		if err != nil {
			fmt.Println(le, err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		typeReq := types.WSBaseReq{}
		err = json.Unmarshal(message, &typeReq)
		if err != nil {
			log.Println("解码失败", err)
			continue
		}

		if typeReq.Type == "AUTH" {
			authReq := types.WSAuthReq{}
			err = json.Unmarshal(message, &authReq)
			if err != nil {
				log.Println("解码失败", err)
				continue
			}

			data, err := logic.JwtAuth(c.svc, authReq.Token)
			if err != nil {
				log.Println("err", err)
				continue
			}

			err = logic.JwtClaimAuth(c.svc, data)
			if err != nil {
				log.Println("err", err)
				continue
			}
			c.JwtId = data.JwtId
			c.UserId = data.UserId

			c.hub.register <- c
			continue
		}
		if typeReq.Type == "WEBRTC" {
			msg := types.WSWebRtcMsg{}
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("解码失败", err)
				continue
			}

			oppoId := msg.OppoId
			msg.OppoId = c.UserId

			//datas, err := json.Marshal(msg)
			//if err != nil {
			//	log.Println("解码失败", err)
			//	continue
			//}

			c.svc.ChannelChat <- &types.ChannelMessage{
				UserId:  oppoId,
				JwtId:   "",
				Content: message,
			}

		}

		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, svc *svc.ServiceContext, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		svc:  svc,
		conn: conn,
		send: make(chan []byte, bufSize),
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
