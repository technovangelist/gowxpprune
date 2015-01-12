package main

import (
	"encoding/xml"
	"fmt"
)

type Category struct {
	XMLName  xml.Name      `xml:"category"`
	TermId   wpstring      `xml:"term_id"`
	NiceName wpstring      `xml:"category_nicename"`
	Parent   wpstring      `xml:"category_parent"`
	Name     wpcdatastring `xml:"cat_name"`
}

func (c Category) CNiceName() wpstring {
	return c.NiceName
}

var categorymap map[wpstring]Category

func (c Category) String() string {
	return fmt.Sprintf("\t ID : %s - NiceName : %s\n", c.TermId, c.NiceName)
}

func create_category_map() {
	categorymap = make(map[wpstring]Category)
	for _, category := range rss.Channel.Categories {
		categorymap[category.TermId] = Category{TermId: category.TermId, NiceName: category.NiceName, Parent: category.Parent, Name: category.Name}
	}
}

func (c *Category) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*c, start)
	return nil //xml.Header + string(out)
}

// func delete_category(category_id int) {
//
//
//
// 	delete(categorymap, category_id)
// 	var NewCategories []Category
// 	for _, i := range categorymap {
// 		NewCategories = append(NewCategories, categorymap[i.TermId])
// 	}
// 	rss.Channel.Categories = NewCategories
// }
