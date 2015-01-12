package main

import (
	"encoding/xml"
	"fmt"
)

type Tag struct {
	XMLName xml.Name      `xml:"tag"`
	TermId  wpstring      `xml:"term_id"`
	Slug    wpstring      `xml:"tag_slug"`
	Name    wpcdatastring `xml:"tag_name"`
}

func (t Tag) String() string {
	return fmt.Sprintf("\t ID : %s - Slug : %s - Name : %s\n", t.TermId, t.Slug, t.Name)
}

func (t *Tag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*t, start)
	return nil //xml.Header + string(out)
}
