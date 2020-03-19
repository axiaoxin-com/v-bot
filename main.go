package main

import (
	"cuitclock/config"
	"log"
)

func main() {
	config.InitConfig()
	log.Println("[INFO] cuitclock inited config.")
	runCronServer()
}
