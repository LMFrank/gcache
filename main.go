package main

import (
	"flag"
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

func createGroup() *sgcache.Group {
	return sgcache.NewGroup("scores", 2<<10, sgcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// 启动缓存服务器
// 创建HTTPPool，添加节点信息，注册到sg中，启动HTTP服务
func startCacheServer(addr string, addrs []string, sg *sgcache.Group) {
	peers := sgcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	sg.RegisterPeers(peers)
	log.Println("sgcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 用来启动一个API服务（端口 9999），与用户进行交互
func startAPIServer(apiAddr string, sg *sgcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := sg.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "SgCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	sg := createGroup()
	if api {
		go startAPIServer(apiAddr, sg)
	}
	startCacheServer(addrMap[port], addrs, sg)
}
