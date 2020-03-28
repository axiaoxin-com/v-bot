mv config.toml config.toml.bak
mv cuitclock cuitclock.bak
tar xzf cuitclock.tar.gz
mv config.toml.build config.toml
supervisorctl restart cuitclock
tail -f /var/log/cuitclock.log
