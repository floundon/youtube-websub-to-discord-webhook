package handler

import (
	"context"
	"fmt"
	"github.com/floundon/youtube-websub-to-discord-webhook/config"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/discordwebhook"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/youtube"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := discordwebhook.SendWithContext(ctx, config.Get().WebHookURL, &discordwebhook.Request{
		Content:   fmt.Sprintf(`
新しいライブ配信・動画が登録されました。
New Live Stream or Video has been added.
Title: %s
YouTube: %s
`, request.Entries[0].Title, request.Entries[0].Link.HRef.String()),
	})

	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}

	c.String(http.StatusOK, "")
}