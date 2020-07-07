tar xzvf cnarea20181031.sql.tar.gz
# dbname cnarea
mysql -uroot -p -e 'create database cnarea'
mysql -uroot -p cnarea < cnarea20181031.sql
