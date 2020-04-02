package progressbar

import (
	"fmt"
	"testing"
	"time"

	"github.com/schollz/progressbar/v2"
)

func TestProgressBar(t *testing.T) {

	theme := progressbar.Theme{Saucer: "+", SaucerHead: "", SaucerPadding: "-", BarStart: "", BarEnd: ""}
	bar := ProgressBar(theme, 10, 100, 25)
	fmt.Println(bar)
}

func TestDaytProgressBar(t *testing.T) {
	bar := DayProgressBar(time.Now())
	fmt.Println(bar)
}

func TestYeartProgressBar(t *testing.T) {
	bar := YearProgressBar(time.Now())
	fmt.Println(bar)
}
