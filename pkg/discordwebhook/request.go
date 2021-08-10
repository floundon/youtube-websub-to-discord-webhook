package discordwebhook

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"unicode/utf8"
)

// MaxContentLength
// Maximum length of content field
// Ref: https://discord.com/developers/docs/resources/webhook#execute-webhook
const MaxContentLength = 2000

const requestContentType = "application/json"

type Request struct {
	Username string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Content string `json:"content"`
}

func (r *Request) Validate() error {
	lengthOfContent := utf8.RuneCountInString(r.Content)
	if lengthOfContent == 0 {
		return xerrors.Errorf("content is empty")
	} else if lengthOfContent > MaxContentLength {
		return xerrors.Errorf("maximum length of contend field exceeded: current=%d, max=%d", lengthOfContent, MaxContentLength)
	}

	return nil
}

func (r *Request) toJSON() ([]byte, error) {
	return json.Marshal(r)
}
