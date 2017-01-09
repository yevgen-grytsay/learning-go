package main

import (
    "path/filepath"
    "os"
    "log"
    "fmt"
    "container/list"
    "bytes"
    "crypto/md5"
    "io"
    "encoding/hex"
    "strings"
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

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func FindDuplicates(l *list.List) {
    similar := make(map[string]*list.List)
    for e := l.Front(); e != nil; e = e.Next() {
        if path, ok := e.Value.(string); ok {
            hash, _ := HashFile(path)
            if _, ok := similar[hash]; !ok {
                similar[hash] = list.New()
            }
            similar[hash].PushBack(path)
        }
    }

    for _, files := range similar {
        if files.Len() > 1 {
            fmt.Println("Duplicated files: ", ListToString(files))
        }
    }
}

//http://www.mrwaggel.be/post/generate-md5-hash-of-a-file/
func HashFile(filePath string) (string, error) {
    fmt.Println("Hashing file", filePath)
    //Initialize variable returnMD5String now in case an error has to be returned
    var returnMD5String string

    //Open the passed argument and check for any error
    file, err := os.Open(filePath)
    if err != nil {
        return returnMD5String, err
    }

    //Tell the program to call the following function when the current function returns
    defer file.Close()

    //Open a new hash interface to write to
    hash := md5.New()

    //Copy the file in the hash interface and check for any error
    if _, err := io.Copy(hash, file); err != nil {
        return returnMD5String, err
    }

    //Get the 16 bytes hash
    //hashInBytes := hash.Sum(nil)[:16]
    hashInBytes := hash.Sum(nil)

    //Convert the bytes to a string
    returnMD5String = hex.EncodeToString(hashInBytes)

    return returnMD5String, nil
}

func main() {
    var cr = Crawler{Sizes: make(map[int64]*list.List)}
    filepath.Walk("/home/yevgen", cr.Visitor)

    for size, items := range cr.Sizes {
        if items.Len() > 1 {
            fmt.Println(strings.Repeat("-", 80))
            fmt.Println("Possible duplicates", size, ListToString(items))
            FindDuplicates(items)
        }
        //fmt.Println(size, " (", items.Len(), "): ", ListToString(items))
    }
}