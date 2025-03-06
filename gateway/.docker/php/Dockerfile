# syntax=docker/dockerfile:1.10-labs
ARG ALPINE_VERSION=3.21
ARG PHP_VERSION=8.4
ARG COMPOSER_VERSION=2.8
ARG PHP_EXT_INSTALLER_VERSION=2.7
ARG ROADRUNNER_VERSION=2024.3

FROM composer:${COMPOSER_VERSION} AS composer_stage
FROM mlocati/php-extension-installer:${PHP_EXT_INSTALLER_VERSION} AS php_ext_installer_stage
FROM ghcr.io/roadrunner-server/roadrunner:${ROADRUNNER_VERSION} AS roadrunner

FROM php:${PHP_VERSION}-cli-alpine${ALPINE_VERSION} AS app_php_pre_composer
RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    apk update \
    && install-php-extensions opcache sockets pcntl pdo_pgsql pdo_mysql redis \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

ARG USER_NAME=app-user
ARG GROUP_NAME=app-user
ARG UID=1000
ARG GID=1000

RUN addgroup -g ${GID} -S ${GROUP_NAME} && adduser -u ${UID} --ingroup ${GROUP_NAME} -S -g ${USER_NAME} ${USER_NAME} \
    && mkdir -p /home/${USER_NAME} /app \
    && chown -R ${USER_NAME}:${GROUP_NAME} /home/${USER_NAME} \
    && chown -R ${USER_NAME}:${GROUP_NAME} /app

WORKDIR /app
COPY --from=roadrunner /usr/bin/rr /app/rr

COPY ./.docker/php/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

USER ${USER_NAME}

FROM app_php_pre_composer AS app_php_composer
USER root

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${USER_NAME}:${GROUP_NAME} /usr/bin/composer
USER ${USER_NAME}


FROM app_php_pre_composer AS app_php
ENV COMPOSER_HOME="/tmp/composer"
COPY --chown=${USER_NAME}:${GROUP_NAME} ./composer.* .
RUN --mount=type=bind,from=composer_stage,source=/usr/bin/composer,target=/usr/bin/composer \
    composer install --no-cache --no-ansi --no-autoloader --no-scripts --no-dev

COPY --exclude=./tests --link --chown=${USER_NAME}:${GROUP_NAME} . .

RUN --mount=type=bind,from=composer_stage,source=/usr/bin/composer,target=/usr/bin/composer \
    set -x \
    && composer dump-autoload -n --optimize


FROM app_php AS app_php_dev
USER root

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${USER_NAME}:${GROUP_NAME} /usr/bin/composer

RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    install-php-extensions xdebug \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

COPY ./.docker/php/config/dev.ini ${PHP_INI_DIR}/conf.d/

USER ${USER_NAME}

FROM app_php AS app_php_prod
USER root

COPY ./.docker/php/config/prod.ini ${PHP_INI_DIR}/conf.d/

USER ${USER_NAME}

FROM app_php AS app_php_test
USER root

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer
RUN chown -R ${USER_NAME}:${GROUP_NAME} /usr/bin/composer
RUN composer install

COPY --link --chown=${USER_NAME}:${GROUP_NAME} ./tests ./tests

COPY ./.docker/php/config/prod.ini ${PHP_INI_DIR}/conf.d/

USER ${USER_NAME}
