package errors

import "net/http"

type Error struct {
	ErrorCode    int    `json:"error_code"`
	ErrorText    string `json:"error_text"`
	DebugMessage string `json:"debug_message"`
}

func NewError(code int, err error, debugMessage string) (int, *Error) {
	return code, &Error{ErrorCode: code, ErrorText: err.Error(), DebugMessage: debugMessage}
}

func NewError400(err error, debugMessage string) (int, *Error) {
	return NewError(http.StatusBadRequest, err, debugMessage)
}

func NewError500(err error, debugMessage string) (int, *Error) {
	return NewError(http.StatusInternalServerError, err, debugMessage)
}
