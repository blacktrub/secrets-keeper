package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"secrets-keeper/pkg/keybuilder"
	"secrets-keeper/pkg/storage"

	"github.com/gin-gonic/gin"
)

const MaxLenghtMessage = 1024
const MaxTTL = 86400
const MinTTL = 60

func validateMessageLenght(msg string) bool {
	return len(msg) <= MaxLenghtMessage
}

func validateTTLSize(ttl int) bool {
	if ttl < MinTTL {
		return false
	}

	if ttl >= MaxTTL {
		return false
	}

	return true
}

func generateLink(c *gin.Context, key string) string {
	return fmt.Sprintf("http://%s/%s", c.Request.Host, key)
}

func writeInternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
}

func writeBadRequest(c *gin.Context, reason string) {
	c.HTML(http.StatusBadRequest, "400.html", gin.H{"reason": reason})
}

func indexView(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{"maxTTL": MaxTTL, "maxMessageLenght": MaxLenghtMessage},
	)
}

func saveMessageView(c *gin.Context, keyBuilder keybuilder.KeyBuilder, keeper keeper.Keeper) {
	message := c.PostForm("message")
	if !validateMessageLenght(message) {
		log.Println("Bad message lenght")
		writeBadRequest(c, "message")
		return
	}

	ttl, err := strconv.Atoi(c.PostForm("ttl"))
	if err != nil {
		ttl = MinTTL
	}

	if !validateTTLSize(ttl) {
		log.Println("Bad ttl")
		writeBadRequest(c, "ttl")
		return
	}

	key, err := keyBuilder.Get()
	if err != nil {
		log.Println("Keybuilder error", err)
		writeInternalError(c)
		return
	}

	err = keeper.Set(key, message, ttl)
	if err != nil {
		log.Println("Keeper error", err)
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "key.html", gin.H{"key": generateLink(c, key)})
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
	router.Run(":8080")
}
