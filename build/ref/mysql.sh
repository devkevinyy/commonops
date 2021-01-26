#!/bin/sh

MYSQL_DATABASE=common_ops
MYSQL_USER=mysql
MYSQL_ROOT_PASSWORD=ops123456ycj

rm -f /etc/my.cnf
cp my.cnf /etc/

ps -ef|grep mysqld|grep -v grep|awk '{print $1}'|xargs kill -9

if [ -d /var/lib/mysql ]; then
    rm -rf /var/lib/mysql/*
else
    mkdir -p /var/lib/mysql
fi

mkdir -p /run/mysqld/ && chmod -R 777 /run/mysqld/
rm -f /run/mysqld/*

echo "[i] MySQL user: $MYSQL_USER"
echo "[i] MySQL user 'mysql' Password: $MYSQL_ROOT_PASSWORD"

echo "" > temp

cat << EOF > temp
    USE mysql;
    FLUSH PRIVILEGES;
    GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY "$MYSQL_ROOT_PASSWORD" WITH GRANT OPTION;
    GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
    UPDATE user SET password=PASSWORD("${MYSQL_ROOT_PASSWORD}") WHERE user='root' AND host='localhost';
EOF

if [ "$MYSQL_DATABASE" != "" ]; then
    echo "[i] Creating database: $MYSQL_DATABASE"
    echo "CREATE DATABASE IF NOT EXISTS \`$MYSQL_DATABASE\` CHARACTER SET utf8 COLLATE utf8_general_ci;" >> temp
fi

if [ "$MYSQL_USER" != "" ]; then
  echo "[i] Creating user: $MYSQL_USER with password $MYSQL_PASSWORD"
  echo "GRANT ALL ON *.* to '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_ROOT_PASSWORD';" >> temp
fi

mysql_install_db

nohup /usr/bin/mysqld &

sleep 2s

echo "mysql server start success!"
echo "mysql data init ..."

cat temp | mysql -uroot

rm -f temp

mysql -uroot -Dcommon_ops < ./common_ops.sql

echo "mysql data init finish!"