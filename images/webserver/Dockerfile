FROM php:7.4-fpm-alpine

COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/bin/supervisord

RUN apk add --no-cache \
      nginx \
      shadow \
      bzip2-dev \
      ca-certificates \
      freetype-dev \
      icu-dev \
      libzip-dev \
      pcre-dev \
      libsodium-dev && \
    apk add --no-cache --virtual build-deps \
      coreutils \
      build-base \
      autoconf \
      automake && \
    docker-php-ext-install -j$(nproc) \
      bcmath \
      intl \
      sockets \
      bz2 \
      sodium \
      zip && \
    update-ca-certificates && \
    ln -s /usr/local/bin/php /usr/bin/php && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log && \
    rm -rf /var/lib/nginx/tmp && \
    ln -sf /tmp /var/lib/nginx/tmp && \
    mkdir /etc/supervisor.d/ /var/tmp/nginx/ && \
    chown -R www-data:www-data /var/lib/nginx /var/tmp/nginx/ && \
    chmod 777 -R /var/tmp/nginx/ && \
    pecl install apcu && \
    docker-php-ext-enable apcu && \
    apk del build-deps && \
    rm -rf /tmp/* && \
    chown -R www-data:www-data /var/www && \
    usermod -u 1000 www-data && \
    apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community gnu-libiconv && \
    echo 'upload_max_filesize = 128M' > /usr/local/etc/php/conf.d/docker.ini && \
    echo 'post_max_size = 128M' >> /usr/local/etc/php/conf.d/docker.ini && \
    echo 'max_execution_time = 30' >> /usr/local/etc/php/conf.d/docker.ini && \
    echo 'date.timezone = "Europe/Berlin"' >> /usr/local/etc/php/conf.d/docker.ini && \
    sed -i -e 's/listen = 9000/listen = \/tmp\/php-fpm.sock/g' /usr/local/etc/php-fpm.d/zz-docker.conf && \
    echo 'listen.mode = 7777' >> /usr/local/etc/php-fpm.d/zz-docker.conf && \
    echo 'user = root' >> /usr/local/etc/php-fpm.d/zz-docker.conf

ENV LD_PRELOAD /usr/lib/preloadable_libiconv.so php

COPY rootfs/etc/ /etc/

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
