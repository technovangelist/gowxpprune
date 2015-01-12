package main

import (
	"encoding/xml"
	"fmt"
)

type Term struct {
	XMLName  xml.Name      `xml:"term"`
	TermId   wpstring      `xml:"term_id"`
	Taxonomy wpstring      `xml:"term_taxonomy"`
	Slug     wpstring      `xml:"term_slug"`
	Name     wpcdatastring `xml:"term_name"`
}

func (t Term) String() string {
	return fmt.Sprintf("\t ID : %s - Taxonomy : %s - Slug : %s - Name : %s\n", t.TermId, t.Taxonomy, t.Slug, t.Name)
}

func (t *Term) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*t, start)
	return nil //xml.Header + string(out)
}
