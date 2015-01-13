package main

import (
	"html"
	"io/ioutil"
	"regexp"
)

func ns_cleanup(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		//Do something
	}

	//add the xml thing at the top
	stringcontent := "<?xml version='1.0' encoding='UTF-8' ?>\n" + string(content)

	//fix the CDATA blocks
	re := regexp.MustCompile(`(&lt;)(!\[CDATA\[)(.*)(\]\])`)
	tweakedstring := re.ReplaceAllStringFunc(stringcontent, func(match string) string {
		matches := re.FindStringSubmatch(match)
		return "<" + matches[2] + html.UnescapeString(matches[3]) + matches[4] + ">"
	})

	//OMG Carriage returns matter?!?!?!
	re = regexp.MustCompile(`(<wp:author>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<\/wp:author>)`)
	tweakedstring = re.ReplaceAllStringFunc(tweakedstring, func(match string) string {
		matches := re.FindStringSubmatch(match)
		return matches[1] + matches[2] + matches[3] + matches[4] + matches[5] + matches[6] + matches[7] + matches[8]
	})
	re = regexp.MustCompile(`(<wp:category>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<\/wp:category>)`)
	tweakedstring = re.ReplaceAllStringFunc(tweakedstring, func(match string) string {
		matches := re.FindStringSubmatch(match)
		return matches[1] + matches[2] + matches[3] + matches[4] + matches[5] + matches[6]
	})
	re = regexp.MustCompile(`(<wp:tag>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<\/wp:tag>)`)
	tweakedstring = re.ReplaceAllStringFunc(tweakedstring, func(match string) string {
		matches := re.FindStringSubmatch(match)
		return matches[1] + matches[2] + matches[3] + matches[4] + matches[5]
	})
	re = regexp.MustCompile(`(<wp:term>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<wp:.*>.*<\/.*>)\s*(<\/wp:term>)`)
	tweakedstring = re.ReplaceAllStringFunc(tweakedstring, func(match string) string {
		matches := re.FindStringSubmatch(match)
		return matches[1] + matches[2] + matches[3] + matches[4] + matches[5] + matches[6]
	})
	// remove the namespace attributes of content and excerpt enconded elements

	re = regexp.MustCompile(`encoded( xmlns="[^"]+")`)

	tweakedstring = re.ReplaceAllString(tweakedstring, "encoded")

	ioutil.WriteFile(filename, []byte(tweakedstring), 0777)

}
