package youtubepubsub

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_UnmarshalXML(t *testing.T) {
	cases := []struct {
		message    string
		input      string
		expectFunc func(t *testing.T, parsedTime time.Time)
	}{
		{
			message: "success(published)",
			input:   "<published>2015-03-06T21:40:57+00:00</published>",
			expectFunc: func(t *testing.T, parsedTime time.Time) {
				assert.Equal(t, parsedTime.Year(), 2015)
				assert.Equal(t, parsedTime.Hour(), 21)
				assert.Equal(t, parsedTime.Minute(), 40)
				assert.Equal(t, parsedTime.Second(), 57)
				assert.Equal(t, parsedTime.Nanosecond(), 0)
			},
		},
		{
			message: "success(updated)",
			input:   "<updated>2015-03-09T19:05:24.552394234+00:00</updated>",
			expectFunc: func(t *testing.T, parsedTime time.Time) {
				assert.Equal(t, parsedTime.Year(), 2015)
				assert.Equal(t, parsedTime.Hour(), 19)
				assert.Equal(t, parsedTime.Second(), 24)
				assert.Equal(t, parsedTime.Nanosecond(), 552394234)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.message, func(t *testing.T) {
			u := &DateTime{}
			err := xml.Unmarshal([]byte(c.input), u)
			assert.Nil(t, err)
			c.expectFunc(t, u.Time)
		})
	}
}
