package main

import (
	"github.com/Jacobbrewer1/repoHandler/config"
	"log"
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
	if err := handleGithub(); err != nil {
		log.Println(err)
		return
	}
}
