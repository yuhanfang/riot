// Launches a ratelimit server on the specified port. See documentation in
// github.com/yuhanfang/riot/ratelimit/service/server for details on server
// interface. See github.com/yuhanfang/riot/ratelimit/service/client for a
// reference client implementation.
//
// Usage example:
// 		ratelimit_server --port=8080
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/yuhanfang/riot/ratelimit/service/server"
)

var port = flag.Int("port", 8080, "server port")

func main() {
	flag.Parse()
	http.Handle("/", server.New())
	log.Println("listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
