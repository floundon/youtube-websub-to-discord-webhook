package youtube

import (
	"encoding/xml"
	"net/url"
)

type URL struct {
	url.URL `xml:"-"`
}

func (u *URL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}
	u.URL = *parsedURL

	return nil
}

func (u *URL) UnmarshalXMLAttr(attr xml.Attr) error {
	parsedURL, err := url.Parse(attr.Value)
	if err != nil {
		return err
	}
	u.URL = *parsedURL

	return nil
}
