package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//单独拎出来hhhh，就只写一个lv1咯

func TestPing(t *testing.T) {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	req, _ := http.NewRequest("GET", "/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// 验证响应
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "pong", resp.Body.String())
}

func TestEcho(t *testing.T) {
	router := gin.Default()

	router.GET("/echo", func(c *gin.Context) {
		body, _ := c.GetQuery("message")
		c.JSON(200, gin.H{
			"message": body,
		})
	})

	req, _ := http.NewRequest("GET", "/echo?message=hello", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// 验证响应
	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, `{"message":"hello"}`, resp.Body.String())
}
