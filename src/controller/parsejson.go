package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strings"
)

func printSyntaxError(path, js string, syntax *json.SyntaxError) {
    start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
    if idx := strings.Index(js[start:], "\n"); idx >= 0 {
        end = start + idx
    }

    line, pos := strings.Count(js[:start], "\n")+1, int(syntax.Offset)-start-1

    fmt.Printf("Fail to load json %s in line %d: %s \n", path, line, syntax.Error())
    fmt.Printf("%s\n%s^", js[start:end], strings.Repeat(" ", pos))
}

func LoadJsonConfig(path string, v interface{}) error {
    file, err := ioutil.ReadFile(path)
    if err != nil {
        println(err.Error())
        return err
    }
    err = json.Unmarshal(file, v)
    if err != nil {
        if serr, ok := err.(*json.SyntaxError); ok {
            printSyntaxError(path, string(file), serr)
            return err
        }
        println(err.Error())
        return err
    }
    return nil
}

