package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig(paths ...string) error {
	processdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(processdir)
	viper.AddConfigPath(workdir)
	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	return viper.ReadInConfig()
}
