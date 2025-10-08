# syntax=docker/dockerfile:1.10-labs
ARG ALPINE_VERSION
ARG PHP_VERSION
ARG COMPOSER_VERSION
ARG PHP_EXT_INSTALLER_VERSION

FROM composer:${COMPOSER_VERSION} AS composer_stage
FROM mlocati/php-extension-installer:${PHP_EXT_INSTALLER_VERSION} AS php_ext_installer_stage

FROM php:${PHP_VERSION}-cli-alpine${ALPINE_VERSION} AS app_php_pre_composer
ARG GROUP_NAME=appuser
ARG USER_NAME=appuser
ARG UID
ARG GID
ARG PROTOC_VERSION

RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    apk update \
    && PHP_VER=$(php -r 'echo PHP_MAJOR_VERSION . PHP_MINOR_VERSION;') \
    && apk add --no-cache "php${PHP_VER}-pecl-grpc" "php${PHP_VER}-pecl-protobuf" "protoc=${PROTOC_VERSION}" "grpc-plugins" \
    && EXT_DIR=$(php-config --extension-dir) \
    && cp "/usr/lib/php${PHP_VER}/modules/grpc.so" "${EXT_DIR}/grpc.so" \
    && cp "/usr/lib/php${PHP_VER}/modules/protobuf.so" "${EXT_DIR}/protobuf.so" \
    && MAKEFLAGS="-j$(nproc)" IPE_SWOOLE_WITHOUT_IOURING=1 install-php-extensions swoole pcntl sockets pdo_pgsql \
    && docker-php-ext-enable grpc protobuf \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

RUN if getent group ${GID} >/dev/null 2>&1; then \
    EXISTING_GROUP=$(getent group ${GID} | cut -d: -f1); \
    adduser -u ${UID} --ingroup $EXISTING_GROUP -S -g "${USER_NAME}" ${USER_NAME}; \
    else \
    addgroup -g ${GID} -S ${GROUP_NAME}; \
    adduser -u ${UID} --ingroup ${GROUP_NAME} -S -g "${USER_NAME}" ${USER_NAME}; \
    fi \
    && mkdir -p /home/${USER_NAME} /app \
    && chown -R ${UID}:${GID} /home/${USER_NAME} /app

WORKDIR /app

COPY ./.docker/gateway/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

FROM app_php_pre_composer AS app_php_composer
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${UID}:${GID} /usr/bin/composer
USER ${UID}

FROM app_php_pre_composer AS app_php_dev
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${UID}:${GID} /usr/bin/composer

RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    MAKEFLAGS="-j$(nproc)" install-php-extensions xdebug \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

COPY ./.docker/gateway/config/dev.ini ${PHP_INI_DIR}/conf.d/

USER ${UID}

FROM app_php_pre_composer AS app_php_test
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"

COPY --chown=${UID}:${GID} ./services/gateway .

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${UID}:${GID} /usr/bin/composer
RUN composer install

COPY ./.docker/gateway/config/php.ini ${PHP_INI_DIR}/conf.d/

USER ${UID}
