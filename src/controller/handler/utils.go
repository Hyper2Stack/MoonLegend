package handler

import (
    "github.com/op/go-logging"
)

const (
    RequestBodyDecodeError = "invalid request body"
    RequestBodyError       = "invalid request body field"
    DuplicateResource      = "duplicate resource"
    InvalidOperation       = "invalid operation"

    TokenExpireTime = 3600000
)

var log, _ = logging.GetLogger("moonlegend")

func encodeUserToken(username string) string {
    // TBD
    return username
}

func decodeUserToken(key string) (string, bool) {
    // TBD
    // username, expiered
    return key, true
}
