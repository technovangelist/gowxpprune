package main

import (
	"encoding/xml"
	"fmt"

	"regexp"
)

type Item struct {
	XMLName       xml.Name      `xml:"item"`
	Title         string        `xml:"title"`
	Link          string        `xml:"link"`
	PubDate       string        `xml:"pubDate"`
	Creator       dcstring      `xml:"creator"`
	Guid          Guid          `xml:"guid"`
	Description   string        `xml:"description"`
	Content       contentstring `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Excerpt       excerptstring `xml:"http://wordpress.org/export/1.2/excerpt/ encoded"`
	PostId        wpstring      `xml:"post_id"`
	PostDate      wpstring      `xml:"post_date"`
	GMTPostDate   wpstring      `xml:"post_date_gmt"`
	CommentStatus wpstring      `xml:"comment_status"`
	PingStatus    wpstring      `xml:"ping_status"`
	PostName      wpstring      `xml:"post_name"`
	Status        wpstring      `xml:"status"`
	PostParent    wpstring      `xml:"post_parent"`
	MenuOrder     wpstring      `xml:"menu_order"`
	PostType      wpstring      `xml:"post_type"`
	PostPassword  wpstring      `xml:"post_password"`
	IsSticky      wpstring      `xml:"is_sticky"`
	AttachmentUrl wpstring      `xml:"attachment_url"`
	PostMetas     []*PostMeta   `xml:"postmeta"`
	Comments      []*Comment    `xml:"comment"`
}

type Guid struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Guid        string `xml:",innerxml"`
}

type PostMeta struct {
	XMLName   xml.Name      `xml:"postmeta"`
	MetaKey   wpstring      `xml:"meta_key"`
	MetaValue wpcdatastring `xml:"meta_value"`
}

type Comment struct {
	XMLName        xml.Name       `xml:"comment"`
	Id             wpstring       `xml:"comment_id"`
	Author         wpstring       `xml:"comment_author"`
	AuthorEmail    wpstring       `xml:"comment_author_email"`
	AuthorUrl      wpstring       `xml:"comment_author_url"`
	AuthorIp       wpstring       `xml:"comment_author_IP"`
	CommentDate    wpstring       `xml:"comment_date"`
	CommentDateGMT wpstring       `xml:"comment_date_gmt"`
	Content        wpstring       `xml:"comment_content"`
	Approved       wpstring       `xml:"comment_approved"`
	Type           wpstring       `xml:"comment_type"`
	Parent         wpstring       `xml:"comment_parent"`
	UserId         wpstring       `xml:"comment_user_id"`
	Meta           []*CommentMeta `xml:"commentmeta"`
}

type CommentMeta struct {
	XMLName   xml.Name `xml:"commentmeta"`
	MetaKey   wpstring `xml:"meta_key"`
	MetaValue wpstring `xml:"meta_value"`
}

func create_item_map() {
	itemmap = make(map[wpstring]*Item)
	postmap = make(map[wpstring]*Item)
	attachmentmap = make(map[wpstring]*Item)
	attachmentmap_by_url = make(map[wpstring]*Item)
	for _, item := range rss.Channel.Items {
		itemmap[item.PostId] = &Item{
			Title:         item.Title,
			Link:          item.Link,
			PubDate:       item.PubDate,
			Creator:       item.Creator,
			Guid:          item.Guid,
			Description:   item.Description,
			Content:       item.Content,
			Excerpt:       item.Excerpt,
			PostId:        item.PostId,
			PostDate:      item.PostDate,
			GMTPostDate:   item.GMTPostDate,
			CommentStatus: item.CommentStatus,
			PingStatus:    item.PingStatus,
			PostName:      item.PostName,
			Status:        item.Status,
			PostParent:    item.PostParent,
			MenuOrder:     item.MenuOrder,
			PostType:      item.PostType,
			PostPassword:  item.PostPassword,
			IsSticky:      item.IsSticky,
			AttachmentUrl: item.AttachmentUrl,
			PostMetas:     item.PostMetas,
			Comments:      item.Comments,
		}
		if item.PostType == "post" {
			postmap[item.PostId] = itemmap[item.PostId]
		} else if item.PostType == "attachment" {
			attachmentmap[item.PostId] = itemmap[item.PostId]
			attachmentmap_by_url[item.AttachmentUrl] = itemmap[item.PostId]
			// fmt.Println(item.AttachmentUrl)
			// fmt.Println(attachmentmap[item.PostId].AttachmentUrl)
			// fmt.Println(attachmentmap_by_url[item.AttachmentUrl].PostId)
		}
	}
}

func (i Item) String() string {
	if i.PostType == "post" {
		return fmt.Sprintf("\nPostId : %s - Title : %s - Link : %s - Type : %s\n", i.PostId, i.Title, i.Link, i.PostType)
	} else {
		return ""
	}

}

// func (i *Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	start.Name.Local = "wp:" + start.Name.Local
// 	e.EncodeElement(*i, start)
// 	return nil //xml.Header + string(out)
// }

func (p *PostMeta) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*p, start)
	return nil //xml.Header + string(out)
}

func (c *Comment) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "wp:" + start.Name.Local
	e.EncodeElement(*c, start)
	return nil //xml.Header + string(out)
}

var itemmap map[wpstring]*Item
var postmap map[wpstring]*Item
var attachmentmap map[wpstring]*Item
var attachmentmap_by_url map[wpstring]*Item

func delete_post_by_id(post_id wpstring) {
	fmt.Printf("\nDeleting post ID: %s - Title: %s", post_id, itemmap[post_id].Title)
	posturls := find_urls_in_post(string(itemmap[post_id].Content), "www.datadoghq.com")

	for _, url := range posturls {
		delete_attachment(wpstring(url))
	}
	delete(itemmap, post_id)

	var NewItems []*Item
	for _, i := range itemmap {
		NewItems = append(NewItems, itemmap[i.PostId])
	}

	rss.Channel.Items = NewItems
	create_maps()
}

func delete_attachment(attachment_url wpstring) {
	// fmt.Printf("\nDeleting attachment: %s\n", attachment_url)
	// log.Println(attachmentmap_by_url)
	if attachment, ok := attachmentmap_by_url[attachment_url]; ok {
		// fmt.Println("Found the attachment")
		delete(itemmap, attachment.PostId)
	} else {
		fmt.Println("Attachment not found")
	}

	var NewItems []*Item
	for _, i := range itemmap {
		NewItems = append(NewItems, itemmap[i.PostId])
	}

	rss.Channel.Items = NewItems
}

func find_urls_in_post(text string, siteurl wpstring) (urls []string) {
	regexstring := fmt.Sprintf(`"(https?\://%s/wp-content/uploads(/\S*)?)"`, siteurl)
	r, err := regexp.Compile(regexstring)

	if err != nil {
		fmt.Printf("bad regex")
	}
	submatched_urls := r.FindAllStringSubmatch(text, -1)
	for _, url := range submatched_urls {
		urls = append(urls, url[1])
	}
	return
}

func cleanup_unused_attachments() {
	cleaned_attachments := 0
	fmt.Println("Looking for unused attachments and removing them. Please wait a few seconds.")
	for _, attachment := range attachmentmap_by_url {
		fmt.Printf(".")
		// fmt.Println("looking for attachment: ", attachment.AttachmentUrl)
		attachment_found := false
		for _, item := range itemmap {

			itemurls := find_urls_in_post(string(item.Content), "www.datadoghq.com")
			for _, url := range itemurls {
				if url == string(attachment.AttachmentUrl) {
					attachment_found = true
				}
			}
		}

		if attachment_found == false {
			delete_attachment(attachment.AttachmentUrl)
			cleaned_attachments += 1
		}

	}
	fmt.Println("\nCleaned Unused Attachments:", cleaned_attachments)
	create_maps()
}
