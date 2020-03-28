package weiboclock

import (
	"testing"
)

func TestRandFont(t *testing.T) {
	f, err := RandFont()
	if err != nil {
		t.Error(err)
	}
	if f == nil {
		t.Error("RandFont return nil font")
	}
}
