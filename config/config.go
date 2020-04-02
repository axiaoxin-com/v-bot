package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig(paths ...string) {
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
	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	viper.SetDefault("weiboclock.timezone", "Asia/Shanghai")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("[FATAL] viper ReadInConfig error", err)
	}
}
