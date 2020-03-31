package weiboclock

import (
	"image/color"
	"os"
	"testing"
	"time"
)

func TestMergeClockPic(t *testing.T) {
	clock, err := os.Open("../assets/images/clock/0.png")
	if err != nil {
		t.Error(err)
	}
	defer clock.Close()
	icon, err := os.Open("../assets/images/clock/icon.jpg")
	//icon, err := os.Open("/Users/ashin/Downloads/x.jpg")
	if err != nil {
		t.Error(err)
	}
	defer icon.Close()
	buffer, err := MergeClockPic(time.Now(), clock, icon, "jpg", color.RGBA{0, 0, 0, 255})
	if err != nil {
		t.Error(err)
	}
	f, err := os.Create("/tmp/new.png")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	f.Write(buffer.Bytes())
	if _, err := os.Stat("/tmp/new.png"); err != nil {
		t.Error(err)
	}
}
