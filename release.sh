#! /bin/bash

env GOOS=linux go build -o cuitclock
tar czvf cuitclock.tar.gz cuitclock config.toml pictures
