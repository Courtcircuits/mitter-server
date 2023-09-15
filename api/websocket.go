package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Courtcircuits/mitter-server/types"
	"github.com/Courtcircuits/mitter-server/util"
	"github.com/gorilla/websocket"
)

type Owner struct {
	id   int
	name string
}

type Connection struct {
	id     int
	Conn   *websocket.Conn
	Hub    *Hub
	Owner  Owner
	authed bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ErrSockReqInvalidFormat = errors.New("the message is not the good format, it must be TextMessage")

func Handler(w http.ResponseWriter, r *http.Request, h *Hub) error {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}
	connection := h.AddConnection(conn)

	log.Println("new connection")

	connection.ReceiveMessages()

	h.RemoveConnection(connection.id)
	return nil
}

func (c *Connection) AuthConn(token string) (Owner, error) {
	clear, err := util.VerifyJWT(token)
	if err != nil {
		log.Println("failed loging in")
		return Owner{}, err
	}
	id, err := strconv.Atoi(clear["id"].(string))
	if err != nil {
		return Owner{}, err
	}
	log.Printf("logged %q \n", clear["name"].(string))
	return Owner{
		id:   id,
		name: clear["name"].(string),
	}, nil
}

func (c *Connection) SendAllMessages() error {
	messages, err := GetServer().store.GetMessages()
	if err != nil {
		return err
	}
	messages_jsonified, err := json.Marshal(messages)

	if err != nil {
		return err
	}

	if err = c.Conn.WriteMessage(websocket.TextMessage, messages_jsonified); err != nil {
		return err
	}
	return nil
}

// send a message to the client
func (c *Connection) SendMessage(msg types.Message) error {
	content, err := json.Marshal([]types.Message{msg})
	if err != nil {
		log.Println(err)
		return err
	}

	if err = c.Conn.WriteMessage(websocket.TextMessage, content); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// receive a message from the client
func (c *Connection) ReceiveMessages() error {
	for {
		messageType, p, err := c.Conn.ReadMessage() //message is a JSON
		if err != nil {
			log.Println(err)
			return err
		}

		if !c.authed {
			owner, err := c.AuthConn(string(p))
			if err != nil {
				log.Println(err)
				return err
			}
			c.Owner = owner
			c.authed = true

			if err = c.SendAllMessages(); err != nil {
				return err
			}
			continue
		}

		if messageType != websocket.TextMessage {
			return ErrSockReqInvalidFormat
		}

		msg, err := serv.store.CreateMessage(string(p), c.Owner.id, c.Owner.name)

		if err != nil {
			log.Println(err)
			return err
		}

		c.Hub.Broadcast(msg)

		if string(p) == "exit" {
			break
		}

	}
	return nil
}
