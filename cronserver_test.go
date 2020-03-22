package main

import (
	"cuitclock/config"
	"testing"
)

func TestDoToll(t *testing.T) {
	config.InitConfig()
	doToll(true)
}
