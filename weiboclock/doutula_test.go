package weiboclock

import (
	"testing"
)

func TestDoutulaSearch(t *testing.T) {
	s, err := DoutulaSearch("1", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
