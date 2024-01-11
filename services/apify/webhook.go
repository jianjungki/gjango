package apify

import (
	"encoding/json"
	"log"
	"tiktok_tools/apperr"

	"github.com/gin-gonic/gin"
)

// NewApifyService create a new apify application service
func NewApifyService() *Service {
	return &Service{}
}

// Service represents the user application service
type Service struct {
}

type ActionWebHook struct {
	UserID    string                 `json:"userId"`
	CreatedAt string                 `json:"createdAt"`
	EventType string                 `json:"eventType"`
	EventData map[string]interface{} `json:"eventData"`
	Resource  map[string]interface{} `json:"resource"`
}

/**
{
    "userId": {{userId}},
    "createdAt": {{createdAt}},
    "eventType": {{eventType}},
    "eventData": {{eventData}},
    "resource": {{resource}}
}
**/

func (s *Service) WebHook(c *gin.Context) {
	var r ActionWebHook
	if err := c.ShouldBindJSON(&r); err != nil {
		log.Fatalf("parse webhook error: %v", err.Error())
		apperr.Response(c, err)
		return
	}
	log.Println("event type: " + r.EventType)
	if r.Resource != nil {
		for header, value := range r.Resource {
			valueStr, _ := json.Marshal(value)
			log.Println("header: " + header + "\tvalue: " + string(valueStr))
		}
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}
