package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zinkt/ginkweb/ginkcache"
)

var db = map[string]string{
	"tom":   "111",
	"jason": "222",
	"zinkt": "777",
}

func main() {
	ginkcache.NewGroup("scores", 2<<10, ginkcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key: ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := ginkcache.NewHTTPPool(addr)
	log.Println("ginkcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
