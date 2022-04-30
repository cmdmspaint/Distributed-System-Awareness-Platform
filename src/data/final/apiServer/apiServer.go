package main

import (
	"Distributed-System-Awareness-Platform/src/data/final/apiServer/heartbeat"
	"Distributed-System-Awareness-Platform/src/data/final/apiServer/locate"
	"Distributed-System-Awareness-Platform/src/data/final/apiServer/objects"
	"Distributed-System-Awareness-Platform/src/data/final/apiServer/versions"
	"Distributed-System-Awareness-Platform/src/data/final/dataServer/temp"
	"log"
	"net/http"
	"os"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
