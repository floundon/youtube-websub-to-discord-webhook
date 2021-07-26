package youtube

import (
	"encoding/xml"
	"time"
)

type DateTime struct {
	time.Time `xml:"-"`
}

func (dt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", value)
	if err != nil {
		return err
	}
	dt.Time = t

	return nil
}
