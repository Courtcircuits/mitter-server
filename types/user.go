package types

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int            `json:"id,omitempty"`
	Name     sql.NullString `json:"name,omitempty"`
	Password sql.NullString `json:"password,omitempty"`
}

func ScanUser(row *sql.Row) (User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func (u *User) Map() gin.H {
	return gin.H{
		"id":   u.ID,
		"name": u.Name,
	}
}
