package BugNetService

import (
	"encoding/base64"
	"io/ioutil"
	"regexp"
)

type ImageSrc struct {
	Ext  string
	Body string
}

// decode body from base64
func (s *ImageSrc) DecodeBody() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Body)
}

// save as file
func (s *ImageSrc) SaveAsFile(fileName string) {
	bytes, err := s.DecodeBody()
	if err == nil {
		ioutil.WriteFile(fileName, bytes, 0666)
	}
}

// get image src from message image
func GetImageSrc(image MessageImage) ImageSrc {
	var ext, body string

	// get ext
	re := regexp.MustCompile(`data:image/(.*);`)
	match := re.FindStringSubmatch(image.ImageToken)
	if len(match) > 1 {
		ext = match[1]
	}

	// get body
	re = regexp.MustCompile(`base64,(.*)"`)
	match = re.FindStringSubmatch(image.ImageToken)
	if len(match) > 1 {
		body = match[1]
	}

	return ImageSrc{Ext: ext, Body: body}
}
