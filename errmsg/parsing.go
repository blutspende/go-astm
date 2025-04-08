package errmsg

import "errors"

const (
	Parsing_MsgInvalidEncoding = "invalid encoding: %s"
)

var (
	Parsing_ErrNotEnoughLines   = errors.New("not enough lines")
	Parsing_ErrInvalidLinebreak = errors.New("invalid line breaking")
)
