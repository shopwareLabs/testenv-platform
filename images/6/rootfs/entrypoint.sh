#!/bin/bash

/usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql &

unset APP_ENV
echo "APP_ENV=prod" > /var/www/shop/.env.local
chown -R 1000 /var/www/shop/.env

while ! mysqladmin ping --silent; do
    sleep 1
done

if [[ -n $APP_URL ]]; then
  mysql -proot shopware -e "UPDATE sales_channel_domain set url = '${APP_URL}' where url = 'http://localhost/shop/public'"
else
  mysql -proot shopware -e "UPDATE sales_channel_domain set url = 'http://${VIRTUAL_HOST}/shop/public' where url = 'http://localhost/shop/public'"
fi

rm -rf /var/www/shop/var/cache/* || true

if [[ -n $SHOPWARE_DEMO_USER_PASSWORD ]]; then
  sudo -E -u www-data /var/www/shop/bin/console frosh:user:change:password demo "$SHOPWARE_DEMO_USER_PASSWORD"
fi

/usr/bin/supervisord -c /etc/supervisord.conf
