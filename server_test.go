package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"secrets-keeper/pkg/keybuilder"
	"secrets-keeper/pkg/storage"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var dummyKeeper = keeper.GetDummyKeeper()

func handleTestRequest(w *httptest.ResponseRecorder, r *http.Request) {
	keyBuilder := keybuilder.GetDummyKeyBuilder()
	router := getRouter(keyBuilder, dummyKeeper)
	router.ServeHTTP(w, r)
}

func TestIndexPage(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("index page is not 200")
	}
}

func TestSaveMessage(t *testing.T) {
	testMessage := "foo"
	postData := strings.NewReader(fmt.Sprintf("message=%s", testMessage))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("save is not 200")
	}

	keyBuilder := keybuilder.GetDummyKeyBuilder()
	key, _ := keyBuilder.Get()
	savedMessage, _ := dummyKeeper.Get(key)
	if savedMessage != testMessage {
		t.Error("message was not saved")
	}

	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)
	if !strings.Contains(string(data), key) {
		t.Error("result page without key")
	}
}

func TestReadMessage(t *testing.T) {
	testMessage := "helloMessage"
	keyBuilder := keybuilder.GetDummyKeyBuilder()
	key, _ := keyBuilder.Get()
	dummyKeeper.Set(key, testMessage)
	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("read message is not 200")
	}

	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)
	if !strings.Contains(string(data), testMessage) {
		t.Error("result page without key")
	}

	_, err := dummyKeeper.Get(key)
	if err == nil {
		t.Error("keeper value must be empty")
	}
}

func TestReadMessageNotFound(t *testing.T) {
	keyBuilder := keybuilder.GetDummyKeyBuilder()
	key, _ := keyBuilder.Get()
	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 404 {
		t.Error("empty message must be 404")
	}
}

func TestOneReader(t *testing.T) {
	dummyKeeper := keeper.GetDummyKeeper()
	testMessage := "helloMessage"
	keyBuilder := keybuilder.GetDummyKeyBuilder()
	key, _ := keyBuilder.Get()
	dummyKeeper.Set(key, testMessage)

	router := getRouter(keyBuilder, dummyKeeper)
	resultChannel := make(chan int, 2)

	go func(key string, c chan int, router gin.Engine) {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		resultChannel <- w.Code
	}(key, resultChannel, *router)

	go func(key string, c chan int, router gin.Engine) {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		resultChannel <- w.Code
	}(key, resultChannel, *router)

	firstCode := <-resultChannel
	secondCode := <-resultChannel

	fmt.Println(firstCode, secondCode)
	if firstCode != 200 {
		t.Error("first must be 200")
	}

	if secondCode != 404 {
		t.Error("first must be 404")
	}
}
