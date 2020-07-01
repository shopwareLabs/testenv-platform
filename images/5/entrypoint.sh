#!/bin/sh

/usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql &

while ! mysqladmin ping --silent; do
    sleep 1
done

mysql shopware -e "UPDATE s_core_shops set host = '${VIRTUAL_HOST}'"

sudo -E -u www-data git clone https://github.com/FriendsOfShopware/FroshAdminer.git /var/www/html/custom/plugins/FroshAdminer
sudo -E -u www-data /var/www/html/bin/console sw:plugin:refresh

sudo -E -u www-data /var/www/html/bin/console sw:plugin:install FroshAdminer --activate

if [[ ! -z $PLUGIN_NAME ]]; then
    sudo -E -u www-data /var/www/html/bin/console sw:plugin:install $PLUGIN_NAME --activate
fi

rm -rf /var/www/html/var/cache/*

/usr/bin/supervisord -c /etc/supervisord.conf
