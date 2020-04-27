package weiboclock

import (
	"log"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestCityAstroInfo(t *testing.T) {
	viper.AddConfigPath("..")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	info, err := CityAstroInfo("深圳", time.Now())
	if err != nil {
		t.Error(err)
	}
	log.Println(info)
}
