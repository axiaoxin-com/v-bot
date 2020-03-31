package weiboclock

import (
	"fmt"
	"testing"
)

func TestProgressBar(t *testing.T) {
	for i := 1; i <= 24; i++ {

		bar := ProgressBar(24, 24, i)
		fmt.Println(bar)
	}
}
