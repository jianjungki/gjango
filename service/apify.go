package service

import (
	"tiktok_tools/services/apify"

	"github.com/gin-gonic/gin"
)

// AuthRouter creates new auth http service
func ApifyRouter(svc *apify.Service, r *gin.Engine) {
	t := ApifyService{svc}
	r.POST("/apify/webhook", t.svc.WebHook)

}

// Auth represents auth http service
type ApifyService struct {
	svc *apify.Service
}
