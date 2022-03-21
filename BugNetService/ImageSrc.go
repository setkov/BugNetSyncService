package BugNetService

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"strings"
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
func (s *ImageSrc) SaveAsFile(fileName string) error {
	bytes, err := s.DecodeBody()
	if err != nil {
		return err
	} else {
		ioutil.WriteFile(fileName, bytes, 0666)
		return nil
	}
}

// get image src from image tag
func GetImageSrc(imageTag string) ImageSrc {
	// generate random image name
	bytes := make([]byte, 16)
	rand.Read(bytes)
	name := "img_" + hex.EncodeToString(bytes)
	// get ext
	ext := getSubstring(imageTag, "data:image/", ";")
	// get body
	body := getSubstring(imageTag, "base64,", "\"")

	return ImageSrc{Name: name, Ext: ext, Body: body}
}

// get first substring between "start" and "end" strings
func getSubstring(str string, start string, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	e += s
	return str[s:e]
}
