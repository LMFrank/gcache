package main

import (
	"fmt"
	"log"
	"net/http"
	"sgcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	sgcache.NewGroup("scores", 2<<10, sgcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:8000"
	peers := sgcache.NewHTTPPool(addr)
	log.Println("sgcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
