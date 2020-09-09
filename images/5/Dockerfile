FROM ghcr.io/shyim/shopware-docker/5/nginx-production:php74

COPY --from=composer /usr/bin/composer /usr/bin/composer

ARG SHOPWARE_DL=https://www.shopware.com/de/Download/redirect/version/sw5/file/install_5.6.8_7b49bfb8ea0d5269b349722157fe324a341ed28e.zip

RUN apk add --no-cache mysql \
    mysql-client \
    sudo \
    git && \
    cd /var/www/html && \
    wget $SHOPWARE_DL && \
    unzip *.zip && \
    rm *.zip && \
    mysql_install_db --datadir=/var/lib/mysql --user=mysql

RUN mkdir /run/mysqld/ && chown -R mysql:mysql /run/mysqld/ && /usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql & sleep 2 && \
    mysql -e "CREATE DATABASE shopware" && \
    mysqladmin --user=root password 'root' && \
    php /var/www/html/recovery/install/index.php --shop-host localhost --db-host localhost --db-socket /run/mysqld/mysqld.sock  --db-user root --db-password root --db-name shopware --shop-locale en_GB --shop-currency EUR --admin-username demo   --admin-password demo --admin-email demo@foo.com --admin-locale en_GB --admin-name demo -n && \
    php /var/www/html/bin/console sw:firstrunwizard:disable && \
    php /var/www/html/bin/console sw:store:download SwagDemoDataEN && \
    php /var/www/html/bin/console sw:plugin:install --activate SwagDemoDataEN && \
    chown -R 1000:1000 /var/www/html

COPY entrypoint.sh /entrypoint.sh

CMD ["/bin/sh", "/entrypoint.sh"]
