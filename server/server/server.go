package server

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"server/db"
	"server/env"
	"strconv"
	"strings"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Panicln("[server]", err)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	fs := static.LocalFile(env.ReactBuildDirectory(), true)

	router.Use(static.Serve("/", fs))

	router.POST("/api/new", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		textBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[server] POST /api/new, error reading request body:", err)
			return
		}

		post, err := db.NewPost(user, string(textBytes))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[server] POST /api/new, error creating new post:", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      post.PostId,
			"text":    post.Text,
			"author":  nil,
			"created": post.TimeCreated,
			"likes":   post.Likes,
			"liked":   false,
		})
	})

	router.POST("/api/edit", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/edit, non-integer postId:", err)
			return
		}

		textBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[server] POST /api/edit, error reading request body:", err)
			return
		}

		post := db.PostById(postId)

		if !post.WasCreatedBy(user) {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/edit, trying to edit the post created by someone else")
			return
		}

		if !post.Edit(string(textBytes)) {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/edit, edit failed")
			return
		}

		c.Data(http.StatusOK, "text/plain", textBytes)
	})

	router.POST("/api/delete", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/delete, non-integer postId:", err)
			return
		}

		post := db.PostById(postId)

		if !post.WasCreatedBy(user) {
			c.String(http.StatusOK, "false")
			log.Println("[server] POST /api/delete, trying to edit the post created by someone else")
			return
		}

		ok = post.Delete()

		c.String(http.StatusOK, strconv.FormatBool(ok))
	})

	router.GET("/api/get", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		var jsonTags []any

		f := func(postId uint64) {
			post := db.PostById(postId)
			realPost, ok := post.Find()
			if !ok {
				return
			}

			other, _ := db.UserByPublicKey(realPost.AuthorId)
			var username any = nil
			if user.PublicIdString() != other.PublicIdString() {
				username, _ = other.FindUsername()
			}

			jsonTags = append(jsonTags, gin.H{
				"id":      realPost.PostId,
				"text":    realPost.Text,
				"author":  username,
				"created": realPost.TimeCreated,
				"likes":   realPost.Likes,
				"liked":   user.Liked(post),
			})
		}

		id := db.PostsCount()

		if id == 0 {
			c.JSON(http.StatusOK, jsonTags)
			return
		}

		for {
			id--

			f(id)

			if id == 0 {
				break
			}
		}

		c.JSON(http.StatusOK, jsonTags)
	})

	router.GET("/api/me", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		username, err := user.FindUsername()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[server] GET /api/me, error finding username:", err)
			return
		}

		c.String(http.StatusOK, username)
	})

	router.POST("/api/like", func(c *gin.Context) {
		user, ok := updateCookie(c)
		if !ok {
			return
		}

		postIdString := c.Request.URL.Query().Get("id")
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/like, non-integer postId:", err)
			return
		}

		valueString := c.Request.URL.Query().Get("value")
		value, err := strconv.ParseBool(valueString)
		if err != nil {
			c.AbortWithStatus(http.StatusExpectationFailed)
			log.Println("[server] POST /api/like, non-boolean value:", err)
			return
		}

		post := db.PostById(postId)

		count, ok := user.Like(post, value)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[server] POST /api/like, like error:", err)
			return
		}

		c.String(http.StatusOK, strconv.FormatUint(count, 10))
	})

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "notFound",
			})
			return
		}

		c.FileFromFS("", fs)
	})
}

func Run() {
	if err := router.Run(":" + env.Port()); err != nil {
		log.Println("[server] listen:", err)
	}
}
