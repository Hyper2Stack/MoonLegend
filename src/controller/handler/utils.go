package handler

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/hex"
    "strconv"
    "time"

    "github.com/op/go-logging"
)

const (
    RequestBodyDecodeError = "invalid request body"
    RequestBodyError       = "invalid request body field"
    InvalidOperation       = "invalid operation"

    TokenEncryptKey        = "0123456789abcdefghijklmnopqrstuv"
    TokenExpireTime        = 3600000
)

var CommonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var log, _ = logging.GetLogger("moonlegend")
var cipherBlock cipher.Block

func Initialize() error {
    b := []byte(TokenEncryptKey)
    var err error
    cipherBlock, err = aes.NewCipher(b)
    return err
}

func encrypt(str string) string {
    cfb := cipher.NewCFBEncrypter(cipherBlock, CommonIV)
    ciphertext := make([]byte, len(str))
    cfb.XORKeyStream(ciphertext, []byte(str))
    return hex.EncodeToString(ciphertext)
}

func decrypt(str string) string {
    bytes, _ := hex.DecodeString(str)
    cfbdec := cipher.NewCFBDecrypter(cipherBlock, CommonIV)
    plaintext := make([]byte, len(bytes))
    cfbdec.XORKeyStream(plaintext, bytes)
    return string(plaintext)
}

func encodeUserToken(username string) string {
    return encrypt(strconv.Itoa(int(time.Now().Unix())) + username)
}

func decodeUserToken(key string) (string, bool) {
    text := decrypt(key)
    t, _ := strconv.Atoi(text[:10])
    u := text[10:]
    if t + TokenExpireTime > int(time.Now().Unix()) {
        return u, true
    }
    return "", false
}
