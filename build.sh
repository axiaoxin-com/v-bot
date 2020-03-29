#! /bin/bash

rm cuitclock.tar.gz
statik -src=./assets
env GOOS=linux go build -x -o cuitclock
cp config.toml config.toml.build
cp ../personal-data/chaojiying_accounts.json .
sed -i '' -e "s/test_username/#test_username/g" config.toml.build
sed -i '' -e "s/test_passwd/#test_passwd/g" config.toml.build
sed -i '' -e 's|accounts_json_path = ".*"|accounts_json_path = "./chaojiying_accounts.json"|g' config.toml.build
tar czvf cuitclock.tar.gz cuitclock config.toml.build chaojiying_accounts.json
rm cuitclock chaojiying_accounts.json config.toml.build
