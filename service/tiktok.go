package service

import (
	"tiktok_tools/services/tiktok"

	"github.com/gin-gonic/gin"
)

// AuthRouter creates new auth http service
func TiktokRouter(svc *tiktok.Service, r *gin.Engine) {
	t := TiktokSevice{svc}
	r.GET("/tiktok/auth", t.svc.AuthLink)
	r.GET("/tiktok/webhook", t.svc.WebHook)

}

// Auth represents auth http service
type TiktokSevice struct {
	svc *tiktok.Service
}
