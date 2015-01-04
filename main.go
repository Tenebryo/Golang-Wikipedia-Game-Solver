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

    //holds all the pages for BFS
    queue := make(chan []string, 100000000)

    //regex to extract exclusively wikipedia links from fetched pages
    reFindLink, err := regexp.Compile("href=\"(/wiki/[^\"/ :#]*)\"")

    if err != nil {
        fmt.Println(err)
        return
    }

    //To prevent undefined behavior of concurrently used maps
    mutex := new(sync.Mutex)

    //list of pages the program has already visited, stored as wikipedia.org suffixes
    visited := make(map[string]bool)

    //link to add suffixes onto
    base_link := "https://en.wikipedia.org"

    //page to start at
    http_start := []string{"/wiki/Raspberry_Pi"}

    //page to find
    link_goal := "/wiki/Egypt"

    queue <- http_start

    done := false

    //Value to keep track of distance from original page
    level := 1

    thread_count := 512

    var solution []string

    for i := 0; i < thread_count; i++ {

    	//Start thread
        fmt.Printf("[%2d] Launching...\n", i)
        go func (id int) {
        	defer func() {
            	fmt.Printf("[%2d] Stopping...\n", id)
        		thread_count--
        	} ()
            fmt.Printf("[%2d] Starting...\n", id)

            for L := range queue {
                if len(L) > level {
                    level++
                    fmt.Println("Level:", level, "Len:", len(queue))
                }
                //Fetch Wikipedia page
                req, err := http.Get(base_link + L[len(L)-1])
                if err == nil {
                	//Extract Links
                    b, _ := ioutil.ReadAll(req.Body)
                    links := reFindLink.FindAllStringSubmatch(string(b), -1)

                    //Add links to queue in batches
                    //mutex lock to prevent concurrent map usage
                    mutex.Lock()
                    if !done {
	                    for t := range links {
	                        _, v := visited[links[t][1]]
	                        if links[t][1] == link_goal {
	                            solution = append(L, links[t][1])
	                            done = true
	                            mutex.Unlock()
	                            return
	                        }
	                        if !v {
	                            visited[links[t][1]] = true
	                            queue <- append(L, links[t][1])
	                        }
	                    }
	                }
                    mutex.Unlock()
                } else {
                    fmt.Println(err)
                }
                if done {
                    return
                }
                time.Sleep(10);
            }

            return;
        }(i)
    }

    //wait for threads to find the solution
    ///TODO: should probably find a better way to do this (maybe a channel?)
    for !done || thread_count > 0 {
        time.Sleep(1)
    }

    fmt.Println(solution)

    fmt.Println("Done!")
}
