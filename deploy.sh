#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"


echo 'Restarting Go...'
cd $DIR/webapp/go/
go build -o isubata
sudo systemctl stop isubata.golang.service
cd $DIR
sudo systemctl restart isubata.golang.service
echo 'Restarted!'

sudo cp $DIR/systemd/* /etc/systemd/system/
sudo systemctl daemon-reload

echo 'Updating config file...'
sudo cp "$DIR/nginx.conf" /etc/nginx/nginx.conf
# sudo cp "$DIR/redis.conf" /etc/redis/redis.conf
sudo cp "$DIR/mysqld.cnf" /etc/mysql/mysql.conf.d/mysqld.cnf
echo 'Updated config file!'

echo 'Restarting middleware services...'
# sudo systemctl restart redis.service
# Save cache
sudo systemctl restart mysql.service
sudo systemctl restart nginx.service
echo 'Restarted!'

echo 'Rotating files'
sudo bash -c 'cp /var/log/nginx/access.log /var/log/nginx/access.log.$(date +%s).$(git rev-parse HEAD) && echo > /var/log/nginx/access.log'
sudo bash -c 'cp /tmp/mysql-slow.sql /tmp/mysql-slow.sql.$(date +%s).$(git rev-parse HEAD) && echo > /tmp/mysql-slow.sql'
echo 'Rotated!'
