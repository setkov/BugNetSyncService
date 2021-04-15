package Common

type errorCategory string

const (
	_error   errorCategory = "Error"
	_warning errorCategory = "Warning"
)

type ErrorWithCategory struct {
	category errorCategory
	text     string
}

func NewError(error string) *ErrorWithCategory {
	return &ErrorWithCategory{
		category: _error,
		text:     error,
	}
}

func NewWarning(error string) *ErrorWithCategory {
	return &ErrorWithCategory{
		category: _warning,
		text:     error,
	}
}

func (e *ErrorWithCategory) Error() string {
	return string(e.category) + ": " + e.text
}
