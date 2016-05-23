package model

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"

    "github.com/op/go-logging"
    "github.com/satori/go.uuid"
)

const (
    hashKey = "0123456789abcdefghijklmnopqrstuv"
)

var log, _ = logging.GetLogger("moonlegend")

////////////////////////////////////////////////////////////////////////////////

type Order struct {
    Columns string
    Sequence string
}

func NewOrder(cols, seq string) *Order {
    return &Order{
        Columns: cols,
        Sequence: seq,
    }
}

////////////////////////////////////////////////////////////////////////////////

type Paging struct {
    Offset int
    Size   int
}

func NewPaging(offset, size int) *Paging {
    return &Paging{
        Offset: offset,
        Size: size,
    }
}

////////////////////////////////////////////////////////////////////////////////

type Condition struct {
    Key   string
    Op    string
    Value interface{}
}

func NewCondition(key, op string, value interface{}) *Condition {
    return &Condition{
        Key: key,
        Op: op,
        Value: value,
    }
}

////////////////////////////////////////////////////////////////////////////////

func GenerateOrderSql(o *Order) string {
    if o != nil {
        return fmt.Sprintf(" ORDER BY %s %s", o.Columns, o.Sequence)
    }

    return ""
}

func GenerateLimitSql(p *Paging) string {
    if p != nil {
        return fmt.Sprintf(" LIMIT %d, %d", p.Offset, p.Size)
    }

    return ""
}

func GenerateWhereSql(cs []*Condition) (string, []interface{}) {
    if len(cs) > 0 {
        ks := make([]string, 0)
        vs := make([]interface{}, 0)
        for _, c := range cs {
            ks = append(ks, fmt.Sprintf("%s%s?", c.Key, c.Op))
            vs = append(vs, c.Value)
        }
        return " WHERE " + strings.Join(ks, " and "), vs
    }

    return "", nil
}

////////////////////////////////////////////////////////////////////////////////

func removeElement(s *[]string, e string) {
    for i := 0; i < len(*s); i++ {
        if (*s)[i] == e {
            *s = append((*s)[:i], (*s)[i+1:]...)
            return
        }
    }
}

func containElement(arr []string, e string) bool {
    for _, s := range arr {
        if s == e {
            return true
        }
    }

    return false
}

func hashPassword(password string) string {
    return hex.EncodeToString(hmac.New(sha256.New, []byte(hashKey)).Sum([]byte(password)))
}

func generateKey() string {
    return uuid.NewV4().String()
}
