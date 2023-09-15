package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Courtcircuits/mitter-server/controllers"
	"github.com/Courtcircuits/mitter-server/storage"
	"github.com/Courtcircuits/mitter-server/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	listenAddr string
	router     *gin.Engine
	store      storage.Database
	hub        *Hub
}

var serv *Server

func GetServer() *Server {
	return serv
}

func NewServer(listenAddr string, store storage.Database) *Server {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{util.Get("CLIENT_URL")}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Origin"}

	r.Use(cors.New(config))

	serv = &Server{
		listenAddr: listenAddr,
		store:      store,
		router:     r,
		hub:        NewHub(),
	}

	return serv
}

func (s *Server) Start() error {

	critical_route := s.router.Group("/")
	critical_route.Use(JWTAuth())

	s.router.POST("/signup", s.Signup)
	s.router.POST("/login", s.Login)
	s.router.GET("ws", s.ChatHandler)
	critical_route.POST("/send", s.WriteMessage)
	critical_route.GET("/messages", s.ReadMessages)

	s.router.Run(s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) ChatHandler(c *gin.Context) {
	err := Handler(c.Writer, c.Copy().Request, s.hub)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
	}
	log.Println("connection closed !")
}

func (s *Server) Signup(c *gin.Context) {
	type Credentials struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatus(400)
	}

	auth_token, err := controllers.SignUpUser(s.store, credentials.Username, credentials.Password)

	if err != nil {
		c.AbortWithStatus(401)
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"token": auth_token,
	})
}

func (s *Server) Login(c *gin.Context) {
	type Credentials struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatus(400)
	}

	auth_token, err := controllers.Authenticate(s.store, credentials.Username, credentials.Password)

	if err != nil {
		c.AbortWithStatus(401)
		if err == sql.ErrNoRows {
			return
		}
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"token": auth_token,
	})
}

func (s *Server) WriteMessage(c *gin.Context) {
	type Message_Request struct {
		Content string `json:"content,omitempty"`
	}

	var message_request Message_Request

	if err := c.BindJSON(&message_request); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	message, err := s.store.CreateMessage(message_request.Content, c.GetInt("id"), c.GetString("name"))

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, message.ToJSON())

}

func (s *Server) ReadMessages(c *gin.Context) {

	i := c.Query("since")

	if i != "" {
		j, err := strconv.ParseInt(i, 10, 64)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		since := time.UnixMilli(j)

		messages, err := s.store.GetMessagesSince(since)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		var list_messages []gin.H

		for _, message := range messages {
			list_messages = append(list_messages, message.ToJSON())
		}
		c.JSON(200, list_messages)
		return

	}

	messages, err := s.store.GetMessages()

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	var list_messages []gin.H

	for _, message := range messages {
		list_messages = append(list_messages, message.ToJSON())
	}
	c.JSON(200, list_messages)

}
