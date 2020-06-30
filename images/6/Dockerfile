FROM shyim/shopware-platform-nginx-production:php74

ARG SHOPWARE_DL=https://www.shopware.com/de/Download/redirect/version/sw6/file/install_6.2.0_1589874223.zip

RUN apk add --no-cache mysql mysql-client git sudo && \
    mkdir -p /var/www/shop/ && \
    cd /var/www/shop/ && \
    wget $SHOPWARE_DL && \
    unzip *.zip && \
    rm *.zip && \
    mysql_install_db --datadir=/var/lib/mysql --user=mysql && \
    echo "pdo_mysql.default_socket=/run/mysqld/mysqld.sock" > /usr/local/etc/php/conf.d/pdo_mysql.ini && \
    echo "mysqli.default_socket=/run/mysqld/mysqld.sock" > /usr/local/etc/php/conf.d/mysqli.ini && \
    echo "memory_limit=1G" > /usr/local/etc/php/conf.d/memory.ini

RUN mkdir /run/mysqld/ && chown -R mysql:mysql /run/mysqld/ && /usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql & sleep 2 && \
    mysql -e "CREATE DATABASE shopware" && \
    mysqladmin --user=root password 'root' && \
    php /var/www/shop/public/recovery/install/index.php --shop-host localhost --db-host localhost --db-socket /run/mysqld/mysqld.sock  --db-user root --db-password root --db-name shopware --shop-locale en-GB --shop-currency EUR --admin-username demo  --admin-password demo --admin-email demo@foo.com --admin-locale en-GB --admin-firstname demo --admin-lastname demo -n && \
    php /var/www/shop/bin/console framework:demodata -p 3000 && \
    php /var/www/shop/bin/console dal:refresh:index && \
    mysql shopware -e "INSERT INTO system_config (id, configuration_key, configuration_value, sales_channel_id, created_at, updated_at) VALUES (X'b3ae4d7111114377af9480c4a0911111', 'core.frw.completedAt', '{\"_value\": \"2019-10-07T10:46:23+00:00\"}', NULL, '2019-10-07 10:46:23.169', NULL);" && \
    php /var/www/shop/bin/console theme:refresh && \
    php /var/www/shop/bin/console theme:change Storefront --all && \
    chown -R 1000:1000 /var/www

COPY patches/ProductAvailableFilter.php /var/www/shop/vendor/shopware/core/Content/Product/SalesChannel/ProductAvailableFilter.php
COPY entrypoint.sh /entrypoint.sh
COPY www.conf /etc/nginx/sites-enabled/www.conf

CMD ["/bin/sh", "/entrypoint.sh"]
