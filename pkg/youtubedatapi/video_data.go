package youtubedatapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type VideoData struct {
	VideoID            string
	IsLive             bool
	ScheduledStartTime time.Time
}

func FetchVideoData(ctx context.Context, apiKey string, videoID string) (*VideoData, error) {
	service, err := youtube.NewService(
		ctx,
		option.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, err
	}

	call := service.Videos.List([]string{"contentDetails", "liveStreamingDetails"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	video := response.Items[0]
	fmt.Printf("%+v\n", video.ContentDetails)
	fmt.Printf("%+v\n", video.LiveStreamingDetails)
	scheduledStartTime, err := time.Parse(time.RFC3339, video.LiveStreamingDetails.ScheduledStartTime)
	if err != nil {
		return nil, err
	}

	return &VideoData{
		VideoID:            videoID,
		IsLive:             true,
		ScheduledStartTime: scheduledStartTime,
	}, nil
}
