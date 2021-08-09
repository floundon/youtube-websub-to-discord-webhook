package handler

import (
	"fmt"
	"github.com/floundon/youtube-websub-to-discord-webhook/config"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/youtube"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscriptionRequest struct {
	VerificationToken string `form:"hub.verify_token"`
	Challenge string `form:"hub.challenge"`
}

func VerifySubscription(c *gin.Context) {
	var request SubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if request.VerificationToken != config.Get().VerificationToken {
		c.String(http.StatusBadRequest, "invalid verification token")
		return
	}

	c.String(http.StatusOK, request.Challenge)
}

func ReceiveNotification(c *gin.Context) {
	var request youtube.Feed
	if err := c.ShouldBindXML(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("%+v\n", request)

	c.String(http.StatusOK, "")
}