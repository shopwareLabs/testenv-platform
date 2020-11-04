#!/bin/sh

/usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql &

unset APP_ENV
echo "APP_ENV=prod" > /var/www/shop/.env
chown -R 1000 /var/www/shop/.env

while ! mysqladmin ping --silent; do
    sleep 1
done

mysql -proot shopware -e "UPDATE sales_channel_domain set url = 'http://${VIRTUAL_HOST}/shop/public' limit 1"

sudo -u www-data git clone https://github.com/FriendsOfShopware/FroshPlatformAdminer.git /var/www/shop/custom/plugins/FroshPlatformAdminer --depth=1
sudo -E -u www-data /var/www/shop/bin/console plugin:refresh
sudo -E -u www-data /var/www/shop/bin/console plugin:install -n --activate FroshPlatformAdminer SaasConnect SwagI18nDutch

/usr/bin/supervisord -c /etc/supervisord.conf
