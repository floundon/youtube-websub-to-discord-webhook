package youtubepubsub

type Link struct {
	Rel  string `xml:"rel,attr"`
	HRef URL    `xml:"href,attr"`
}
