package main

import (
	"encoding/xml"
	"fmt"
)

type Author struct {
	XMLName           xml.Name
	AuthorID          wpstring      `xml:"author_id"`
	AuthorLogin       wpstring      `xml:"author_login"`
	AuthorEmail       wpstring      `xml:"author_email"`
	AuthorDisplayName wpcdatastring `xml:"author_display_name"`
	AuthorFirst       wpcdatastring `xml:"author_first_name"`
	AuthorLast        wpcdatastring `xml:"author_last_name"`
}

func (a Author) String() string {
	return fmt.Sprintf("\t ID : %s - FirstName : %s - LastName : %s - UserName : %s \n", a.AuthorID, a.AuthorFirst, a.AuthorLast, a.AuthorLogin)
}

func (a *Author) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*a, start)
	return nil //xml.Header + string(out)
}
