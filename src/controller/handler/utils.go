package handler

import (
    "path/filepath"
)

const (
    RequestBodyDecodeError = "invalid request body"
    RequestBodyError       = "invalid request body field"
    DuplicateResource      = "duplicate resource"
    InvalidOperation       = "invalid operation"

    TokenExpireTime = 3600000
)

func encodeUserToken(username string) string {
    // TBD
    return username
}

func decodeUserToken(key string) (string, bool) {
    // TBD
    // username, expiered
    return key, true
}

func hashPassword(password string) string {
    // TBD
    return password
}

func absCleanPath(path string) string {
	absPath, _ := filepath.Abs(path)
	return filepath.Clean(absPath)
}
