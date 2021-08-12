package youtubepubsub

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"feed"`

	Link []Link `xml:"link"`

	Title string `xml:"title"`

	Entries []Entry `xml:"entry"`

	Updated DateTime `xml:"updated"`
}
