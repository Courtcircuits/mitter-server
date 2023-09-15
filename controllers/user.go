package controllers

import (
	"strconv"
	"time"

	"github.com/Courtcircuits/mitter-server/storage"
	"github.com/Courtcircuits/mitter-server/types"
	"github.com/Courtcircuits/mitter-server/util"
)

func SignUpUser(s storage.Database, name string, password string) (string, error) {
	user, err := s.CreateUser(name, password)
	if err != nil {
		return "", err
	}
	auth_token := create_jwt_user(user)
	return auth_token, nil
}

func Authenticate(s storage.Database, name string, password string) (string, error) {
	user, err := s.FindUser(name, password)
	if err != nil {
		return "", err
	}
	auth_token := create_jwt_user(user)
	return auth_token, nil
}

func create_jwt_user(user types.User) string {
	return util.GenJWT(time.Now().Add(time.Hour*24*7), map[string]any{
		"id":   strconv.Itoa(user.ID),
		"name": user.Name.String,
	})
}
