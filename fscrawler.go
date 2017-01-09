package main

import (
    "path/filepath"
    "os"
    "log"
    "fmt"
    "container/list"
    "bytes"
)


type Crawler struct {
    Sizes map[int64]*list.List
}

func (cr Crawler) Visitor(path string, info os.FileInfo, err error) error {
    if err != nil {
        log.Print(err)
        return nil
    }
    if !info.IsDir() {
        cr.Push(path, info)
        //fmt.Printf("%s (%d)\n", path, info.Size())
    }

    return nil
}

func (cr Crawler) Push(path string, info os.FileInfo) {
    size := info.Size()
    if _, ok := cr.Sizes[size]; !ok {
        cr.Sizes[size] = list.New()
    }
    cr.Sizes[size].PushBack(path)
}

func ListToString(l *list.List) string {
    bs := bytes.NewBufferString("")
    for e := l.Front(); e != nil; e = e.Next() {
        if str, ok := e.Value.(string); ok {
            bs.WriteString(str)
        }

    }
    return bs.String()
}

func main() {
    var cr = Crawler{Sizes: make(map[int64]*list.List)}
    filepath.Walk("/home/yevgen/Music", cr.Visitor)
    for size, items := range cr.Sizes {
        fmt.Println(size, " (", items.Len(), "): ", ListToString(items))
    }
}