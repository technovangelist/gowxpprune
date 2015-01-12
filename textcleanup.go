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

	// remove the namespace attributes of content and excerpt enconded elements

	re = regexp.MustCompile(`encoded( xmlns="[^"]+")`)

	tweakedstring = re.ReplaceAllString(tweakedstring, "encoded")

	ioutil.WriteFile(filename, []byte(tweakedstring), 0777)

}
