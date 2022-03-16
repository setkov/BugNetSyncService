package BugNetService

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"regexp"
)

type ImageSrc struct {
	Name string
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
	var name, ext, body string

	// generate random image name
	bytes := make([]byte, 16)
	rand.Read(bytes)
	name = "img_" + hex.EncodeToString(bytes)

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

	return ImageSrc{Name: name, Ext: ext, Body: body}
}
