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
    cd /var/www/shop/ && \
    php /var/www/shop/public/recovery/install/index.php --shop-host localhost --db-host localhost --db-socket /run/mysqld/mysqld.sock  --db-user root --db-password root --db-name shopware --shop-locale en-GB --shop-currency EUR --admin-username demo  --admin-password demo --admin-email demo@foo.com --admin-locale en-GB --admin-firstname demo --admin-lastname demo -n && \

    wget https://www.shopware.com/de/Download/redirect/version/sw6/file/update_6.2.2_1592398977.zip && \
    unzip -o update_6.2.2_1592398977.zip && \
    php /var/www/shop/public/recovery/update/index.php -n && \
    rm -rf update-assets && \

    wget https://www.shopware.com/de/Download/redirect/version/sw6/file/update_6.2.3_1594641397.zip && \
    unzip -o update_6.2.3_1594641397.zip && \
    php /var/www/shop/public/recovery/update/index.php -n && \
    rm -rf update-assets && \

    wget https://www.shopware.com/de/Download/redirect/version/sw6/file/update_v6.3.0.0_83857c5fb3cde7fc77392d843c45d95b661c36a8.zip && \
    unzip -o update_v6.3.0.0_83857c5fb3cde7fc77392d843c45d95b661c36a8.zip && \
    php /var/www/shop/public/recovery/update/index.php -n && \
    rm -rf update-assets && \

    wget https://www.shopware.com/de/Download/redirect/version/sw6/file/update_v6.3.0.1_8c1c2d2a8fa041a354e77305901ad4597717e756.zip && \
    unzip -o update_v6.3.0.1_8c1c2d2a8fa041a354e77305901ad4597717e756.zip && \
    php /var/www/shop/public/recovery/update/index.php -n && \
    rm -rf update-assets && \
    rm update*.zip && \

    wget https://releases.shopware.com/sw6/update_v6.3.1.0_363abad22982da639888b36d84d52730e8dc05e6.zip && \
    unzip -o update_v6.3.1.0_363abad22982da639888b36d84d52730e8dc05e6.zip && \
    php /var/www/shop/public/recovery/update/index.php -n && \
    rm -rf update-assets && \
    rm update*.zip && \

    php /var/www/shop/bin/console framework:demodata -p 3000 && \
    php /var/www/shop/bin/console dal:refresh:index && \
    mysql shopware -e "INSERT INTO system_config (id, configuration_key, configuration_value, sales_channel_id, created_at, updated_at) VALUES (X'b3ae4d7111114377af9480c4a0911111', 'core.frw.completedAt', '{\"_value\": \"2019-10-07T10:46:23+00:00\"}', NULL, '2019-10-07 10:46:23.169', NULL);" && \
    php /var/www/shop/bin/console theme:refresh && \
    php /var/www/shop/bin/console theme:change Storefront --all && \
    chown -R 1000:1000 /var/www

COPY fix-language.php /
RUN /usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql & sleep 2 && \
    php /fix-language.php

COPY entrypoint.sh /entrypoint.sh
COPY SwagI18nDutch /var/www/shop/custom/plugins/SwagI18nDutch
COPY www.conf /etc/nginx/sites-enabled/www.conf

CMD ["/bin/sh", "/entrypoint.sh"]
