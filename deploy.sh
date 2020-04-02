mv config.toml config.toml.bak
mv v-bot v-bot.bak
tar xzf v-bot.tar.gz
mv config.toml.build config.toml
supervisorctl restart v-bot
tail -f /var/log/v-bot.log
