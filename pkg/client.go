package pkg

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	username string
	roomID   string
	room     *Room
	userID   string
	isClosed bool
	mu       sync.Mutex
}

func (c *Client) writePump() {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.mu.Lock()
			if c.isClosed {
				c.mu.Unlock()
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.isClosed = true
				c.mu.Unlock()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.isClosed = true
				c.mu.Unlock()
				return
			}

			w.Write(message)

		case <-ticker.C:
			c.mu.Lock()
			if c.isClosed {
				c.mu.Unlock()
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			// write ping to websocket connection to keep alive. this is just some
			// basic ws stuff bro, pff do you even code?
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.isClosed = true
				c.mu.Lock()
				return
			}

			c.mu.Unlock()
		}
	}

}

func (c *Client) readPump() {

}
