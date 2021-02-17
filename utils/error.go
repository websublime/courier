package utils

type Exception struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Exception string `json:"exception"`
}

func NewException(exception string, status int, message string) *Exception {
	e := &Exception{
		Exception: exception,
		Code:      status,
		Message:   message,
	}

	return e
}

func (e *Exception) Error() string {
	return e.Message
}
