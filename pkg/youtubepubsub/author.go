package youtubepubsub

import (
	"encoding/xml"
)

type Author struct {
	XMLName xml.Name `xml:"author"`

	Name string `xml:"name"`

	URI URL `xml:"uri"`
}
