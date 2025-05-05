package lineseparator

var LF = string(byte(0x0A))
var CR = string(byte(0x0D))
var LFCR = string([]byte{byte(0x0A), byte(0x0D)})
var CRLF = string([]byte{byte(0x0D), byte(0x0A)})
