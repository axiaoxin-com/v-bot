#! /bin/bash

rm v-bot.tar.gz
statik -src=./assets
env GOOS=linux go build -x -o v-bot
cp config.toml config.toml.build
cp ../../axiaoxin/personal-data/chaojiying_accounts.json .
sed -i '' -e "s/test_username/#test_username/g" config.toml.build
sed -i '' -e "s/test_passwd/#test_passwd/g" config.toml.build
sed -i '' -e 's|accounts_json_path = ".*"|accounts_json_path = "./chaojiying_accounts.json"|g' config.toml.build
tar czvf v-bot.tar.gz v-bot config.toml.build chaojiying_accounts.json
rm v-bot chaojiying_accounts.json config.toml.build
