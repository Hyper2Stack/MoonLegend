package model

import (
    "testing"
)

func index(s string, arr []string) int {
    for i, e := range arr {
        if e == s {
            return i
        }
    }

    return -1
}

func Test_TopologicSort_1(t *testing.T) {
    nodes := []string{"A", "AB", "BCD", "B", "C", "BC", "D", "E", "F"}
    edges := []*Edge{
        &Edge{From: "A", To: "AB"},
        &Edge{From: "B", To: "AB"},
        &Edge{From: "B", To: "BC"},
        &Edge{From: "C", To: "BC"},
        &Edge{From: "BC", To: "BCD"},
        &Edge{From: "D", To: "BCD"},
    }


    result, err := TopologicSort(nodes, edges)
    if err != nil {
        t.FailNow()
    }

    errCon1 := index("A", result) > index("AB", result)
    errCon2 := index("B", result) > index("AB", result)
    errCon3 := index("B", result) > index("BC", result)
    errCon4 := index("C", result) > index("BC", result)
    errCon5 := index("BC", result) > index("BCD", result)
    errCon6 := index("D", result) > index("BCD", result)

    if errCon1 || errCon2 || errCon3 || errCon4 || errCon5 || errCon6 {
        t.FailNow()
    }
}

func Test_TopologicSort_2(t *testing.T) {
    nodes := []string{"A", "AB", "BCD", "B", "C", "BC", "D", "E", "F"}
    edges := []*Edge{
        // cycle graph
        &Edge{From: "A", To: "AB"},
        &Edge{From: "B", To: "AB"},
        &Edge{From: "B", To: "BC"},
        &Edge{From: "C", To: "BC"},
        &Edge{From: "BC", To: "BCD"},
        &Edge{From: "D", To: "BCD"},
        &Edge{From: "BCD", To: "B"},
    }


    _, err := TopologicSort(nodes, edges)
    if err == nil {
        t.FailNow()
    }
}
