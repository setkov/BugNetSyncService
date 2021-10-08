package WebUI

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// Convert text to HTML
func TextToHtml(str string) template.HTML {
	return template.HTML(str)
}

// Unescape
func Unescape(str string) string {
	var unescape string
	unescape = strings.Replace(str, "&nbsp;", " ", -1)
	unescape = strings.Replace(unescape, "&quot;", "'", -1)
	return unescape
}

// Remove html tags
func HtmlToText(str string) string {
	var text string
	var skip bool

	for _, v := range Unescape(str) {
		s := string(v)
		if s == "<" {
			skip = true
		} else if s == ">" {
			skip = false
		} else if !skip {
			text += s
		}
	}
	return text
}

// Trim text
func TrimText(str string, length int) string {
	var trim string
	num := 0

	for _, v := range HtmlToText(str) {
		trim += string(v)
		num++
		if num > length {
			break
		}
	}
	return trim
}

// DateTime format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%v", t.Format("2006-01-02 15:04"))
}
