package BugNetService

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
	var startPosition int
	var tag string

	for i, char := range html {
		if char == '<' {
			startPosition = i
			tag = string(char)
		} else if len(tag) > 0 {
			tag += string(char)
			if char == '>' {
				if len(tag) > 3 && tag[:4] == "<img" {
					src := GetImageSrc(tag)
					image := MessageImage{StartPosition: startPosition, StopPosition: i, ImageTag: tag, ImageSrc: src}
					images.Images = append(images.Images, image)
				}
				tag = ""
			}
		}
	}

	return images
}
