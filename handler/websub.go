package handler

import (
	"context"
	"fmt"
	"github.com/floundon/youtube-websub-to-discord-webhook/config"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/discordwebhook"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/youtube"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"log"
	"net/http"
	"time"
)

type WebSubHandler struct {
	YouTubeVideoDataTable dynamo.Table
}

type SubscriptionRequest struct {
	VerificationToken string `form:"hub.verify_token"`
	Challenge         string `form:"hub.challenge"`
}

type youTubeVideoData struct {
	VideoID string `dynamo:"VideoID"`
}

func (*WebSubHandler) VerifySubscription(c *gin.Context) {
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

func (h *WebSubHandler) ReceiveNotification(c *gin.Context) {
	var request youtube.Feed
	if err := c.ShouldBindXML(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(request.Entries) == 0 {
		c.String(http.StatusOK, "no entries to notify")
		return
	}

	fmt.Printf("%+v\n", request)

	for _, entry := range request.Entries {
		func() {
			var videoData youTubeVideoData
			err := h.YouTubeVideoDataTable.Get("VideoID", entry.YouTubeVideoID).One(&videoData)
			if err == nil {
				log.Printf("notification already sent for Video ID: %s\n", entry.YouTubeVideoID)
				return
			}

			videoData.VideoID = entry.YouTubeVideoID

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err = discordwebhook.SendWithContext(ctx, config.Get().WebHookURL, &discordwebhook.Request{
				Content: fmt.Sprintf(`
新しいライブ配信・動画が登録されました。
New Live Stream or Video has been added.
YouTube: %s
`, entry.Link.HRef.String()),
			})

			if err != nil {
				log.Printf("error: %s", err.Error())
			} else {
				err = h.YouTubeVideoDataTable.Put(&videoData).Run()
				if err != nil {
					log.Printf("dynamodb put error: %s", err.Error())
				}
			}
		}()
	}

	c.String(http.StatusOK, "ok")
}
