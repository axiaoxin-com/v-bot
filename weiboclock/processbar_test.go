package weiboclock

import (
	"fmt"
	"testing"
)

func TestProgressBar(t *testing.T) {
	for i := 0; i < 24; i++ {

		fmt.Println(i)
		bar := ProgressBar(24, 24, i)
		if i == 0 {
			bar = ProgressBar(24, 24, 24)
		}
		fmt.Println(bar)
	}
}
