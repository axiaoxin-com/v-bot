package weiboclock

import (
	"log"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetAstroInfo(t *testing.T) {
	info := GetAstroInfo(104.066541, 30.572269, time.Now())
	log.Println(info)
}
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
