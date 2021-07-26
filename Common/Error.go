package Common

type errorCategory string

const (
	Error   errorCategory = "Error"
	Warning errorCategory = "Warning"
)

type ErrorWithCategory struct {
	category errorCategory
	text     string
}

func NewError(error string) *ErrorWithCategory {
	return &ErrorWithCategory{
		category: Error,
		text:     error,
	}
}

func NewWarning(error string) *ErrorWithCategory {
	return &ErrorWithCategory{
		category: Warning,
		text:     error,
	}
}

func (e *ErrorWithCategory) Error() string {
	return string(e.category) + ": " + e.text
}

func (e *ErrorWithCategory) Category() errorCategory {
	return e.category
}
