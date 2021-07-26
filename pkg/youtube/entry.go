package youtube

type Entry struct {
	ID string `xml:"id"`

	// YouTubeVideoID Video ID
	YouTubeVideoID string `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`

	// YouTubeChannelID Channel ID
	YouTubeChannelID string `xml:"http://www.youtube.com/xml/schemas/2015 channelId"`

	// Title Video Title
	Title string `xml:"title"`

	// Link Video Link
	Link Link `xml:"link"`

	// Author Author information
	Author Author `xml:"author"`

	// PublishedDate Published date
	PublishedDate DateTime `xml:"published"`

	// UpdatedDate Updated date
	UpdatedDate DateTime `xml:"updated"`
}
