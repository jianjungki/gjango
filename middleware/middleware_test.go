package middleware_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	mw "tiktok_tools/middleware"
)

func TestAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw.Add(r, gin.Logger())
}
