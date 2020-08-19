#!/bin/sh

/usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql &

while ! mysqladmin ping --silent; do
    sleep 1
done

mysql -proot shopware -e "UPDATE sales_channel_domain set url = 'http://${VIRTUAL_HOST}/shop/public' limit 1"

sudo -u www-data git clone https://github.com/FriendsOfShopware/FroshPlatformAdminer.git /var/www/shop/custom/plugins/FroshPlatformAdminer --depth=1
sudo -u www-data git clone https://github.com/shopware/app-system.git /var/www/shop/custom/plugins/SaasConnect --depth=1

sudo -E -u www-data /var/www/shop/bin/console plugin:refresh
sudo -E -u www-data /var/www/shop/bin/console plugin:install --activate FroshPlatformAdminer SaasConnect SwagI18nDutch

sudo -E -u www-data echo 'MAILER_URL=smtp://mail:25' >> /var/www/shop/.env

/usr/bin/supervisord -c /etc/supervisord.conf
