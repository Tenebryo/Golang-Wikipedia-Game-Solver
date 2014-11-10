package main

import (
    "fmt"
    "time"
    "sync"
    "regexp"
    "net/http"
    "io/ioutil"
)

func main() {
    thread_count := 0
    queue := make(chan []string, 100000000)

    reFindLink, err := regexp.Compile("href=\"(/wiki/[^\"/ :#]*)\"")

    if err != nil {
        fmt.Println(err)
        return
    }

    mutex := new(sync.Mutex)

    visited := make(map[string]bool)

    base_link := "https://en.wikipedia.org"

    http_start := []string{"/wiki/Raspberry_Pi"}

    link_goal := "/wiki/Tool"

    queue <- http_start

    done := false

    level := 1

    var solution []string

    for i := 0; i < 512:; i++ {
        thread_count++
        fmt.Printf("[%2d] Launching...\n", i)
        go func (id int) {
            
            fmt.Printf("[%2d] Starting...\n", id)

            for L := range queue {
                if len(L) > level {
                    level++
                    fmt.Println("Level:", level, "Len:", len(queue))
                }
                req, err := http.Get(base_link + L[len(L)-1])
                if err == nil {
                    b, _ := ioutil.ReadAll(req.Body)
                    links := reFindLink.FindAllStringSubmatch(string(b), -1)
                    for t := range links {
                        _, v := visited[links[t][1]]
                        if links[t][1] == link_goal {
                            solution = append(L, links[t][1])
                            done = true
                            return
                        }
                        if !v {
                            mutex.Lock()
                            visited[links[t][1]] = true
                            mutex.Unlock()
                            queue <- append(L, links[t][1])
                        }
                    }
                } else {
                    fmt.Println(err)
                }
                if done {
                    return
                }
                time.Sleep(10);
            }

            thread_count--
            return;
        }(i)
    }

    for !done {
        time.Sleep(10)
    }

    fmt.Println(solution)

    fmt.Println("Done!")
}
