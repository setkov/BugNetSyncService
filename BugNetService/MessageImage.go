package BugNetService

import (
	"strings"
)

type MessageImage struct {
	StartPosition int
	StopPosition  int
	ImageTag      string
	ImageSrc      ImageSrc
}

type MessageImages struct {
	Images []MessageImage
}

// get message imges from html string
func GetMessageImages(html string) MessageImages {
	var images MessageImages
	var i, start, stop int

	for {
		i = strings.Index(html[start:], "<img")
		if i == -1 {
			break
		}
		start += i

		i = strings.Index(html[start:], ">")
		if i == -1 {
			break
		}
		stop = start + i + 1

		tag := html[start:stop]
		src := GetImageSrc(tag)
		if src.Body != "" {
			image := MessageImage{StartPosition: start, StopPosition: stop, ImageTag: tag, ImageSrc: src}
			images.Images = append(images.Images, image)
		}

		start = stop + 1
		if start >= len(html) {
			break
		}
	}

	return images
}
