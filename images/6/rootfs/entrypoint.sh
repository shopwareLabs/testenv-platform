#!/bin/sh

/usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql &

unset APP_ENV
echo "APP_ENV=prod" > /var/www/shop/.env
chown -R 1000 /var/www/shop/.env

while ! mysqladmin ping --silent; do
    sleep 1
done

mysql -proot shopware -e "UPDATE sales_channel_domain set url = 'http://${VIRTUAL_HOST}/shop/public' where url = 'http://localhost/shop/public'"

sudo -u www-data git clone https://github.com/FriendsOfShopware/FroshPlatformAdminer.git /var/www/shop/custom/plugins/FroshPlatformAdminer --depth=1
rm -rf /var/www/shop/var/cache/* || true


if sudo -E -u www-data /var/www/shop/bin/console store:download -p SwagI18nDutch; then
    sudo -E -u www-data /var/www/shop/bin/console plugin:refresh
    sudo -E -u www-data /var/www/shop/bin/console plugin:install -n --activate SwagI18nDutch
else
    sudo -E -u www-data /var/www/shop/bin/console store:download -p SwagLanguagePack
    sudo -E -u www-data /var/www/shop/bin/console plugin:refresh
    sudo -E -u www-data /var/www/shop/bin/console plugin:install -n --activate SwagLanguagePack
fi


sudo -E -u www-data /var/www/shop/bin/console plugin:install -n --activate FroshPlatformAdminer SwagDemoProducts

/usr/bin/supervisord -c /etc/supervisord.conf
