package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/db"
	"server/env"
)

func updateCookie(c *gin.Context) (user *db.User, ok bool) {
	ok = true
	secretKey, err := c.Cookie("secretKey")
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			log.Panicln("[server.updateCookie] error getting cookie:", err)
		}
		user, secretKey, err = db.NewUser()
		if err != nil {
			log.Panicln("[server.updateCookie] error creating user:", err)
		}

		c.SetCookie("secretKey", secretKey, 0, "/", env.Domain(), false, true)
		return
	}

	user, err = db.UserBySecretKey(secretKey)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, false
	}

	return
}
