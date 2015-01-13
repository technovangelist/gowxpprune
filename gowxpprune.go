package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type RSS struct {
	XMLName   xml.Name    `xml:"rss"`
	Version   string      `xml:"version,attr"`
	Excerptns xmlnsstring `xml:"excerpt,attr"`
	Contentns xmlnsstring `xml:"content,attr"`
	Wfw       xmlnsstring `xml:"wfw,attr"`
	Dc        xmlnsstring `xml:"dc,attr"`
	Wp        xmlnsstring `xml:"wp,attr"`
	Channel   Channel     `xml:"channel"`
}

type wpstring string
type wpcdatastring string
type dcstring string
type excerptstring string
type contentstring string
type wfwstring string
type xmlnsstring string

var starting_status string

var site_url string

var rss RSS

type Channel struct {
	XMLName     xml.Name    `xml:"channel"`
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Description string      `xml:"description"`
	PubDate     string      `xml:"pubDate"`
	Language    string      `xml:"language"`
	WXRVersion  wpstring    `xml:"wxr_version"`
	SiteUrl     wpstring    `xml:"base_site_url"`
	BlogUrl     wpstring    `xml:"base_blog_url"`
	Authors     []*Author   `xml:"author"`
	Categories  []*Category `xml:"category"`
	Tags        []*Tag      `xml:"tag"`
	Terms       []*Term     `xml:"term"`
	Generator   string      `xml:"generator"`
	Items       []*Item     `xml:"item"`
}

// type content struct {
// 	Data string `xml:",chardata"`
// }

func ParseWXP(filename string) (rss RSS) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	// totalauthors := 0

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "rss" {
				// 			var channel Channel
				var rss RSS
				decoder.DecodeElement(&rss, &startElement)

				return rss
			}
		}
	}
	return

}

func write_new_wxp(rss RSS) {
	fmt.Println(starting_status)
	fmt.Printf("\nYou now have %d categories and %d items including %d attachments and %d posts\n", len(categorymap), len(itemmap), len(attachmentmap), len(postmap))

	filename := "newfile.xml"
	file, _ := os.Create(filename)

	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "  ")
	if err := enc.Encode(rss); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	file.Close()
}

func (s xmlnsstring) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	return xml.Attr{
			Name:  xml.Name{Local: "\n\txmlns:" + name.Local},
			Value: fmt.Sprintf("%s", s)},
		nil
}

func (s wpstring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(string(s), start)
	return nil
}

func (s wpcdatastring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	newstring := `<![CDATA[` + string(s) + "]]"
	e.EncodeElement(newstring, start)
	return nil
}

func (s dcstring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "dc:" + start.Name.Local
	e.EncodeElement(string(s), start)
	return nil
}

func (s excerptstring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "excerpt:" + start.Name.Local
	e.EncodeElement(string(s), start)
	return nil
}

func (s contentstring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "content:" + start.Name.Local
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: ""}, Value: ""}}
	newstring := `<![CDATA[` + string(s) + "]]"
	e.EncodeElement(newstring, start)
	return nil
}

func (s wfwstring) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wfw:" + start.Name.Local
	e.EncodeElement(string(s), start)
	return nil
}
func create_maps() {
	create_item_map()
	create_category_map()
}

func show_main_prompt() {
	if len(starting_status) == 0 {
		starting_status = fmt.Sprintf("\nYou are starting with %d categories and %d items including %d attachments and %d posts\n", len(categorymap), len(itemmap), len(attachmentmap), len(postmap))
		fmt.Println(starting_status)
	} else {
		fmt.Printf("\nYou currently have %d categories and %d items including %d attachments and %d posts\n", len(categorymap), len(itemmap), len(attachmentmap), len(postmap))
	}
	fmt.Printf("(D)elete, (Q)uit\n")

}

func main() {
	loopit := true
	// reader := bufio.NewReader(os.Stdin)
	var input_file_name string
	var output_file_name string
	flag.StringVar(&input_file_name, "i", "", "Input file name")
	flag.StringVar(&output_file_name, "o", "newfile.xml", "Output file name")
	flag.Parse()
	rss = ParseWXP(input_file_name)
	create_maps()
	site_url = strings.TrimPrefix(strings.TrimPrefix(string(rss.Channel.Link), "http://"), "https://")

	fmt.Println("Welcome to the WXP Pruner")
	show_main_prompt()
	cleanup_unused_attachments()

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter text: ")
	// text, _ := reader.ReadString('\n')

	for loopit {
		show_main_prompt()
		var ln string
		fmt.Scanf("%s", &ln)

		switch ln {
		case "q":
			loopit = false
		case "d":
			delete_command()
		default:
			loopit = true
		}
	}

	// rss.Channel.Categories = append(rss.Channel.Categories, Category{TermId: 555, NiceName: "fudge"})
	// delete_category(13)
	write_new_wxp(rss)
	ns_cleanup(output_file_name)
}

func delete_command() {
	var ln string
	fmt.Println("What do you want to delete? (P)ost")
	fmt.Scanf("%s", &ln)
	switch ln {
	case "p":
		delete_post_command()
	}
}

func delete_post_command() {
	var ln string
	fmt.Println("What post do you want to delete. Enter a number or range (x-y)")
	// var keys []wpstring
	// for k := range postmap {
	// 	keys = append(keys, k)
	// }
	// sort.Ints(keys)
	for _, post := range postmap {
		fmt.Println(string(post.PostId) + ": " + post.Title)
	}
	fmt.Scanf("%s", &ln)
	if strings.Index(ln, "-") == -1 {
		delete_post_by_id(wpstring(ln))
	} else {
		post_range := strings.Split(ln, "-")
		first, _ := strconv.ParseInt(post_range[0], 10, 0)
		last, _ := strconv.ParseInt(post_range[1], 10, 0)
		fmt.Printf("Delete posts %d through %d, %d posts\n", first, last, last-first)
		for postid := first; postid < last; postid++ {
			spostid := strconv.FormatInt(postid, 10)
			if _, ok := postmap[wpstring(spostid)]; ok {
				fmt.Println("Found a post")
				delete_post_by_id(wpstring(spostid))
				//do something here
			}
		}
	}

	// func Index(s, sep string) int

}
