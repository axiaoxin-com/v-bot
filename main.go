package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func init() {
	processdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("[FATAL] get processdir error", err)
	}
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("[FATAL] get workdir error", err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(processdir)
	viper.AddConfigPath(workdir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("[FATAL] viper ReadInConfig error", err)
	}
}

func main() {
	log.Println("running cron server...")
	runCronServer()
}
