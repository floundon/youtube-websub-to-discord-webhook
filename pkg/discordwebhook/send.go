package discordwebhook

import (
	"bytes"
	"context"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

func Send(webhookURL string, request *Request) error {
	if err := request.Validate(); err != nil {
		return err
	}

	jsonData, err := request.toJSON()
	if err != nil {
		return err
	}

	_, err = http.Post(webhookURL, requestContentType, bytes.NewBuffer(jsonData))
	return err
}

func SendWithContext(ctx context.Context, webhookURL string, request *Request) error {
	if err := request.Validate(); err != nil {
		return err
	}

	jsonData, err := request.toJSON()
	if err != nil {
		return err
	}

	_, err = ctxhttp.Post(ctx, http.DefaultClient, webhookURL, requestContentType, bytes.NewBuffer(jsonData))
	return err
}