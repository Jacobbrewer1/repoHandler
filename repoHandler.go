package main

import (
	"github.com/Jacobbrewer1/repoHandler/config"
	"log"
	"sync"
)

func init() {
	log.Println("Initializing logging")
	//log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.Println("Logging initialized")
}

func main() {
	if err := config.ReadConfig(); err != nil {
		log.Println(err)
		return
	}
	var w sync.WaitGroup
	w.Add(1)
	go handleGithub(&w)
	w.Wait()
}
