FROM php:7.4-fpm-alpine

ENV TZ=Europe/Berlin \
    APP_URL=http://localhost/shop/public \
    APP_ENV=prod \
    APP_SECRET=440dec3766de53010c5ccf6231c182acfc90bd25cff82e771245f736fd276518 \
    INSTANCE_ID=10612e3916e153dd3447850e944a03fabe89440970295447a30a75b151bd844e \
    DATABASE_URL=mysql://root:root@localhost/shopware \
    MAILER_URL=smtp://mail:25 \
    SHOPWARE_ES_HOSTS= \
    SHOPWARE_ES_ENABLED=0 \
    SHOPWARE_ES_INDEXING_ENABLED=0 \
    SHOPWARE_ES_INDEX_PREFIX= \
    COMPOSER_HOME=/tmp/composer \
    SHOPWARE_HTTP_CACHE_ENABLED=1 \
    SHOPWARE_HTTP_DEFAULT_TTL=7200 \
    BLUE_GREEN_DEPLOYMENT=0 \
    INSTALL_LOCALE=nl-NL \
    INSTALL_CURRENCY=EUR \
    FPM_PM=dynamic \
    FPM_PM_MAX_CHILDREN=5 \
    FPM_PM_START_SERVERS=2 \
    FPM_PM_MIN_SPARE_SERVERS=1 \
    FPM_PM_MAX_SPARE_SERVERS=3 \
    PHP_MAX_UPLOAD_SIZE=128m \
    PHP_MAX_EXECUTION_TIME=300 \
    PHP_MEMORY_LIMIT=512m \
    LD_PRELOAD="/usr/lib/preloadable_libiconv.so php"

COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/bin/supervisord
COPY --from=composer /usr/bin/composer /usr/local/bin/composer
COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/

RUN apk add --no-cache \
      nginx \
      shadow \
      unzip \
      wget \
      sudo \
      bash \
      patch \
      jq \
      mysql=10.5.9-r0 \
      mysql-client=10.5.9-r0 \
      git && \
    install-php-extensions bcmath gd intl mysqli pdo_mysql sockets bz2 gmp soap zip gmp ffi redis opcache calendar && \
    ln -s /usr/local/bin/php /usr/bin/php && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log && \
    rm -rf /var/lib/nginx/tmp && \
    ln -sf /tmp /var/lib/nginx/tmp && \
    mkdir -p /var/tmp/nginx/ || true && \
    chown -R www-data:www-data /var/lib/nginx /var/tmp/nginx/ && \
    chmod 777 -R /var/tmp/nginx/ && \
    rm -rf /tmp/* && \
    mkdir -p /var/www/shop/ && \
    chown -R www-data:www-data /var/www && \
    usermod -u 1000 www-data && \
    mysql_install_db --datadir=/var/lib/mysql --user=mysql && \
    mkdir /run/mysqld/ && chown -R mysql:mysql /run/mysqld/ && \
    apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/v3.13/community gnu-libiconv

ARG SHOPWARE_DL=https://www.shopware.com/de/Download/redirect/version/sw6/file/install_6.2.0_1589874223.zip

COPY rootfs/fix-install.php /fix-install.php
COPY rootfs/usr /usr

RUN /usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql & sleep 2 && \
    mysqladmin --user=root password 'root' && \
    cd /var/www/shop/ && \
    wget -qq $SHOPWARE_DL && \
    unzip -qq *.zip && \
    rm *.zip && \
    php bin/console system:install --create-database --force && \
    php bin/console system:generate-jwt-secret && \
    php bin/console user:create "demo" --admin --password="demo" -n && \
    php bin/console sales-channel:create:storefront --name=Storefront --url="http://localhost/shop/public" && \
    php bin/console theme:change --all Storefront && \
    php bin/console framework:demodata -p 3000 && \
    php bin/console dal:refresh:index && \
    php /fix-install.php && \
    mysql shopware -e "INSERT INTO system_config (id, configuration_key, configuration_value, sales_channel_id, created_at, updated_at) VALUES (X'b3ae4d7111114377af9480c4a0911111', 'core.frw.completedAt', '{\"_value\": \"2019-10-07T10:46:23+00:00\"}', NULL, '2019-10-07 10:46:23.169', NULL);" && \
    chown -R 1000 /var/www/shop

COPY rootfs /
    
EXPOSE 80

COPY --chown=1000:1000 SwagDemoProducts /var/www/shop/custom/plugins/SwagDemoProducts
STOPSIGNAL SIGKILL

ENTRYPOINT ["/entrypoint.sh"]

HEALTHCHECK --timeout=10s CMD curl --silent --fail http://127.0.0.1:80/shop/public/admin
