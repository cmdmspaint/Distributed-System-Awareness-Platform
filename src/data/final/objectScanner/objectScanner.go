package main

import (
	"Distributed-System-Awareness-Platform/src/data/final/apiServer/objects"
	"Distributed-System-Awareness-Platform/src/data/src/lib/es"
	"Distributed-System-Awareness-Platform/src/data/src/lib/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
	}
}

func verify(hash string) {
	log.Println("verify", hash)
	size, e := es.SearchHashSize(hash)
	if e != nil {
		log.Println(e)
		return
	}
	stream, e := objects.GetStream(hash, size)
	defer stream.Close()
	if e != nil {
		log.Println(e)
		return
	}
	d := utils.CalculateHash(stream)
	if d != hash {
		log.Printf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}

}
