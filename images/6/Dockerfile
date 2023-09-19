FROM php:8.1-fpm-alpine

ENV TZ=Europe/Berlin \
    APP_URL=http://localhost/shop/public \
    APP_ENV=prod \
    APP_SECRET=440dec3766de53010c5ccf6231c182acfc90bd25cff82e771245f736fd276518 \
    INSTANCE_ID=10612e3916e153dd3447850e944a03fabe89440970295447a30a75b151bd844e \
    DATABASE_URL=mysql://root:root@localhost/shopware \
    MAILER_URL=smtp://mail:1025 \
    SHOPWARE_ES_HOSTS=es \
    SHOPWARE_ES_ENABLED=0 \
    SHOPWARE_ES_INDEXING_ENABLED=0 \
    SHOPWARE_ES_INDEX_PREFIX=shop \
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
    IPE_GD_WITHOUTAVIF=1 \
    LD_PRELOAD="/usr/lib/preloadable_libiconv.so php"

COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/bin/supervisord
COPY --from=composer/composer:2-bin /composer /usr/local/bin/composer
COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/
COPY --from=ghcr.io/shyim/gnu-libiconv:v3.14 /gnu-libiconv-1.15-r3.apk /gnu-libiconv-1.15-r3.apk
COPY --from=ghcr.io/friendsofshopware/shopware-cli /usr/local/bin/shopware-cli /usr/local/bin/shopware-cli

RUN apk add --no-cache \
      nginx \
      shadow \
      unzip \
      wget \
      sudo \
      bash \
      patch \
      jq \
      mariadb \
      mariadb-client \
      git && \
    apk add --no-cache --allow-untrusted /gnu-libiconv-1.15-r3.apk && rm /gnu-libiconv-1.15-r3.apk && \
    install-php-extensions bcmath gd intl mysqli pdo_mysql sockets bz2 soap zip gmp ffi redis opcache calendar && \
    ln -s /usr/local/bin/php /usr/bin/php && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log && \
    rm -rf /var/lib/nginx/tmp && \
    ln -sf /tmp /var/lib/nginx/tmp && \
    mkdir -p /var/tmp/nginx/ && \
    chown -R www-data:www-data /var/lib/nginx /var/tmp/nginx/ && \
    chmod 777 -R /var/tmp/nginx/ && \
    rm -rf /tmp/* && \
    mkdir -p /var/www/ && \
    chown -R www-data:www-data /var/www && \
    usermod -u 1000 www-data && \
    mysql_install_db --datadir=/var/lib/mysql --user=mysql && \
    mkdir /run/mysqld/ && chown -R mysql:mysql /run/mysqld/ && \
    apk add --no-cache icu-data-full

ARG SHOPWARE_VERSION=6.5.5.1

COPY rootfs/usr /usr

RUN /usr/bin/mysqld --basedir=/usr --datadir=/var/lib/mysql --plugin-dir=/usr/lib/mysql/plugin --user=mysql & sleep 2 && \
    mysqladmin --user=root password 'root' && \
    cd /var/www/ && \
    shopware-cli project create shop $SHOPWARE_VERSION && \
    cd shop && \
    php bin/console system:install --create-database --force --shop-locale nl-NL && \
    php bin/console system:generate-jwt-secret || true && \
    php bin/console user:create "demo" --admin --password="demodemo" -n && \
    php bin/console sales-channel:create:storefront --name=Storefront --url="http://localhost/shop/public" && \
    php bin/console theme:change --all Storefront && \
    composer req --dev shopware/dev-tools --no-scripts && \
    php -dmemory_limit=2G bin/console -e prod framework:demodata --products 3000 && \
    php bin/console -e prod dal:refresh:index && \
    mysql shopware -e "INSERT INTO system_config (id, configuration_key, configuration_value, sales_channel_id, created_at, updated_at) VALUES (X'b3ae4d7111114377af9480c4a0911111', 'core.frw.completedAt', '{\"_value\": \"2019-10-07T10:46:23+00:00\"}', NULL, '2019-10-07 10:46:23.169', NULL);" && \
    php bin/console system:config:set core.metrics.shareUsageData false --json && \
    composer req "frosh/tools:*" "frosh/adminer-platform:*" --no-scripts  && \
    composer remove --dev shopware/dev-tools --no-scripts && \
    composer install --no-dev --no-scripts && \
    php bin/console store:download -p SwagLanguagePack && \
    git clone https://github.com/shopwareLabs/SwagTestEnvironment.git /var/www/shop/custom/plugins/SwagTestEnvironment --depth=1 && \
    php bin/console plugin:refresh && \
    php bin/console plugin:install -n --activate SwagLanguagePack FroshPlatformAdminer SwagTestEnvironment FroshTools

COPY rootfs /

RUN chown -R 1000 /var/www/shop /tmp/composer

EXPOSE 80
WORKDIR /var/www/shop

STOPSIGNAL SIGKILL

ENTRYPOINT ["/entrypoint.sh"]

HEALTHCHECK --timeout=10s CMD curl --silent --fail http://127.0.0.1:80/shop/public/admin
