package BugNetService

// type Src struct {
// 	ext  string
// 	body string
// }

// parse src from string
// func (s *Src) ParseFromString(str string) error {

// 	return Common.NewError("Error on parse src from string")
// }

// // decode body from base64
// func (s *Src) DecodeFromBase64() ([]byte, error) {
// 	return base64.StdEncoding.DecodeString(s.body)
// }

type MessageImage struct {
	StartPosition int
	StopPosition  int
	Token         string
	Atl           string
	Src           string
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
				if len(token) > 3 &&  token[:4] == "<img" {
					image := MessageImage{StartPosition: startPosition, StopPosition: i, Token: token}
					images.Images = append(images.Images, image)
				}
				token = ""
			}
		}
	}

	return images
}

