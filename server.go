package main

import (
	"fmt"
	"net/http"
	"strconv"

	"secrets-keeper/pkg/keybuilder"
	"secrets-keeper/pkg/storage"

	"github.com/gin-gonic/gin"
)

var MESSAGE_MAX_LENGHT = 1024
var MAX_TTL = 86400

func validateMessageLenght(msg string) bool {
	return len(msg) <= MESSAGE_MAX_LENGHT
}

func validateTTLSize(ttl int) bool {
	return ttl <= MAX_TTL
}

func writeInternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
}

func writeBadRequest(c *gin.Context, reason string) {
	c.HTML(http.StatusBadRequest, "400.html", gin.H{"reason": reason})
}

func indexView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func saveMessageView(c *gin.Context, keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper) {
	message := c.PostForm("message")
	if !validateMessageLenght(message) {
		writeBadRequest(c, "message")
		return
	}

	ttl, err := strconv.Atoi(c.PostForm("ttl"))
	if err != nil {
		ttl = 0
	}

	if !validateTTLSize(ttl) {
		writeBadRequest(c, "ttl")
		return
	}

	key, err := keyBuilder.Get()
	if err != nil {
		writeInternalError(c)
		return
	}

	err = keeper.Set(key, message, ttl)
	if err != nil {
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", c.Request.Host, key)})
}

func readMessageView(c *gin.Context, keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper) {
	key := c.Param("key")
	msg, err := keeper.Get(key)
	if err != nil {
		if err.Error() == "not_found" {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "message.html", gin.H{"message": msg})
}

func buildHandler(fn func(c *gin.Context, keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper), keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, keyBuilder, keeper)
	}
}

func getRouter(keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/400.html",
		"templates/500.html",
	)
	router.GET("/", indexView)
	router.POST("/", buildHandler(saveMessageView, keyBuilder, keeper))
	router.GET("/:key", buildHandler(readMessageView, keyBuilder, keeper))
	return router
}

func main() {
	keyBuilder := keybuilder.UUIDKeyBuilder{}
	keeper := keeper.GetRedisKeeper()
	router := getRouter(keyBuilder, keeper)
	router.Run("localhost:8080")
}
