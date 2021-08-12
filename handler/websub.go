package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/youtubedatapi"

	"github.com/floundon/youtube-websub-to-discord-webhook/config"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/discordwebhook"
	"github.com/floundon/youtube-websub-to-discord-webhook/pkg/youtubepubsub"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

type WebSubHandler struct {
	YouTubeVideoDataTable dynamo.Table
}

type SubscriptionRequest struct {
	VerificationToken string `form:"hub.verify_token"`
	Challenge         string `form:"hub.challenge"`
}

type youTubeVideoData struct {
	VideoID     string    `dynamo:"VideoID"`
	ScheduledAt time.Time `dynamo:"ScheduledAt"`
}

const defaultTimeLayout = "2006-01-02 15:04"

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
	var request youtubepubsub.Feed
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
			alreadySent := false
			err := h.YouTubeVideoDataTable.Get("VideoID", entry.YouTubeVideoID).One(&videoData)
			if err == nil {
				alreadySent = true
			}

			videoID := entry.YouTubeVideoID
			videoData.VideoID = videoID

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer fetchCancel()

			liveData, err := youtubedatapi.FetchVideoData(fetchCtx, config.Get().YouTubeAPIKey, videoID)
			if err != nil {
				log.Printf("fetch video data error: %s", err.Error())
				return
			}

			discordTemplate := `
新しいライブ配信・動画が登録されました。 / New Live Stream or Video has been added.
開始時刻 / Scheduled at: %s(JST) / %s(UTC)
%s
`

			// 既に通知が送られていて、開始時刻が一緒の場合は通知しない
			if alreadySent {
				if videoData.ScheduledAt == liveData.ScheduledStartTime {
					log.Printf("notification is already sent: videoID=%s", videoID)
					return
				} else {
					discordTemplate = `
ライブ配信の開始時間が変更されました。 / Live Stream schedule has changed.
開始時刻 / Scheduled at: %s(JST) / %s(UTC)
%s
`
				}
			}
			videoData.ScheduledAt = liveData.ScheduledStartTime

			utc := time.FixedZone("UTC", 0)
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)

			utcStartTime := liveData.ScheduledStartTime.In(utc)
			jstStartTime := liveData.ScheduledStartTime.In(jst)

			postCtx, postCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer postCancel()

			err = discordwebhook.SendWithContext(postCtx, config.Get().WebHookURL, &discordwebhook.Request{
				Content: fmt.Sprintf(
					discordTemplate,
					jstStartTime.Format(defaultTimeLayout),
					utcStartTime.Format(defaultTimeLayout),
					entry.Link.HRef.String(),
				),
			})

			if err != nil {
				log.Printf("discord webhook post error: %s", err.Error())
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
