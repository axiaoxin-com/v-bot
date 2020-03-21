#! /bin/bash

rm cuitclock.tar.gz
statik -src=./assets
env GOOS=linux go build -x -o cuitclock
tar czvf cuitclock.tar.gz cuitclock config.toml
rm cuitclock
