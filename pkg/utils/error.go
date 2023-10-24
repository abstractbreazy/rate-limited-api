package utils

type HTTPError struct {
	Code    int
	Message string
}

func NewHTTPError(code int, message string) (err HTTPError) {
	err = HTTPError{}
	err.Code = code
	err.Message = message
	return
}
