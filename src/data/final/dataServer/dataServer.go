package main

import (
	"Distributed-System-Awareness-Platform/src/data/final/dataServer/heartbeat"
	"Distributed-System-Awareness-Platform/src/data/final/dataServer/locate"
	"Distributed-System-Awareness-Platform/src/data/final/dataServer/objects"
	"Distributed-System-Awareness-Platform/src/data/final/dataServer/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
