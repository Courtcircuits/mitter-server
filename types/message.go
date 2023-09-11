package types

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID         int    `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Name_owner string `json:"name___owner,omitempty"`
}

func ScanMessage(row *sql.Row, name_creator string) (Message, error) {
	type message_db struct {
		id        int
		content   sql.NullString
		timestamp sql.NullTime
		id_owner  int
	}
	var message message_db
	err := row.Scan(&message.id, &message.content, &message.timestamp, &message.id_owner)

	real_message := Message{
		ID:         message.id,
		Content:    message.content.String,
		Timestamp:  strconv.FormatInt(message.timestamp.Time.UnixMilli(), 10),
		Name_owner: name_creator,
	}

	return real_message, err
}

func (m *Message) ToJSON() gin.H {
	return gin.H{
		"id":        m.ID,
		"content":   m.Content,
		"timestamp": m.Timestamp,
		"author":    m.Name_owner,
	}
}
