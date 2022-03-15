package BugNetService

type MessageImage struct {
	StartPosition int
	StopPosition  int
	ImageToken    string
}

type MessageImages struct {
	Images []MessageImage
}

// get message imges from html string
func GetMessageImages(html string) MessageImages {
	var images MessageImages
	var startPosition int
	var token string

	for i, char := range html {
		if char == '<' {
			startPosition = i
			token = string(char)
		} else if len(token) > 0 {
			token += string(char)
			if char == '>' {
				if len(token) > 3 && token[:4] == "<img" {
					image := MessageImage{StartPosition: startPosition, StopPosition: i, ImageToken: token}
					images.Images = append(images.Images, image)
				}
				token = ""
			}
		}
	}

	return images
}
