package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscriptionRequest struct {
	VerifyToken string `form:"hub.verify_token"`
	Challenge string `form:"hub.challenge"`
}

func VerifySubscription(c *gin.Context) {
	var request SubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func ReceiveNotification(c *gin.Context) {

}