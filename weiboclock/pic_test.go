package weiboclock

import (
	"os"
	"testing"
)

func TestDoutulaSearch(t *testing.T) {
	s, err := DoutulaSearch("1", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}

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
	buffer, err := MergeClockPic(clock, icon, "jpg")
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

func TestRandFont(t *testing.T) {
	f, err := RandFont()
	if err != nil {
		t.Error(err)
	}
	if f == nil {
		t.Error("RandFont return nil font")
	}
}
