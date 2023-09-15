package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Courtcircuits/mitter-server/util"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth_header := ctx.GetHeader("Authorization")

		if auth_header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Auth header is empty",
			})
			return
		}

		parts := strings.Split(auth_header, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Auth header must provide a bearer token",
			})
		}

		parsed_token, err := util.VerifyJWT(parts[1])

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}

		id, err := strconv.Atoi(parsed_token["id"].(string))

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "failed cast of id",
			})
		}

		log.Printf("Authenticated %d, %q\n", id, parsed_token["name"])

		ctx.Set("id", id)
		ctx.Set("name", parsed_token["name"])

		ctx.Next()
	}
}
