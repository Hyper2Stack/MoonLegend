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

////////////////////////////////////////////////////////////////////////////////

type Edge struct {
    From string
    To   string
}

func removeEdge(from, to string, edges *[]*Edge) {
    for i := 0; i < len(*edges); i++ {
        if (*edges)[i].From == from && (*edges)[i].To == to {
            *edges = append((*edges)[:i], (*edges)[i+1:]...)
            return
        }
    }
}

func hasParent(node string, edges []*Edge) bool {
    for _, e := range edges {
        if node == e.To {
            return true
        }
    }

    return false
}

func children(node string, edges []*Edge) []string {
    result := make([]string, 0)
    for _, e := range edges {
        if node == e.From {
            result = append(result, e.To)
        }
    }

    return result
}

func TopologicSort(nodes []string, edges []*Edge) ([]string, error) {
    l := make([]string, 0)
    s := make([]string, 0)

    for _, node := range nodes {
        if !hasParent(node, edges) {
            s = append(s, node)
        }
    }

    for {
        if len(s) == 0 {
            break
        }

        n := s[0]
        s = s[1:]
        l = append(l, n)
        for _, m := range children(n, edges) {
            removeEdge(n, m, &edges)
            if !hasParent(m, edges) {
                s = append(s, m)
            }
        }
    }

    if len(edges) != 0 {
        return nil, fmt.Errorf("graph has at least one cycle")
    }

    return l, nil
}
