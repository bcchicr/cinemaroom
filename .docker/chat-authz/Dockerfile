# syntax=docker/dockerfile:1.10-labs
ARG ALPINE_VERSION
ARG PHP_VERSION
ARG COMPOSER_VERSION
ARG PHP_EXT_INSTALLER_VERSION
ARG ROADRUNNER_VERSION
ARG GOLANG_VERSION

FROM composer:${COMPOSER_VERSION} AS composer_stage
FROM mlocati/php-extension-installer:${PHP_EXT_INSTALLER_VERSION} AS php_ext_installer_stage
FROM ghcr.io/roadrunner-server/roadrunner:${ROADRUNNER_VERSION} AS roadrunner

FROM php:${PHP_VERSION}-cli-alpine${ALPINE_VERSION} AS app_php_pre_composer
ARG GROUP_NAME=appuser
ARG USER_NAME=appuser
ARG UID
ARG GID
ARG PROTOC_VERSION
ARG ROADRUNNER_VERSION

RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    apk update \
    && PHP_VER=$(php -r 'echo PHP_MAJOR_VERSION . PHP_MINOR_VERSION;') \
    && apk add --no-cache "php${PHP_VER}-pecl-grpc" "php${PHP_VER}-pecl-protobuf" "protoc=${PROTOC_VERSION}" \
    && EXT_DIR=$(php-config --extension-dir) \
    && cp "/usr/lib/php${PHP_VER}/modules/grpc.so" "${EXT_DIR}/grpc.so" \
    && cp "/usr/lib/php${PHP_VER}/modules/protobuf.so" "${EXT_DIR}/protobuf.so" \
    && MAKEFLAGS="-j$(nproc)" install-php-extensions sockets pcntl pdo_pgsql \
    && docker-php-ext-enable grpc protobuf \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

RUN arch=$(uname -m) && \
    if [ "$arch" = "x86_64" ]; then arch="amd64"; elif [ "$arch" = "aarch64" ]; then arch="arm64"; else echo "Unsupported architecture: $arch" && exit 1; fi \
    && apk add --no-cache wget tar ca-certificates \
    && wget -q "https://github.com/roadrunner-server/roadrunner/releases/download/v${ROADRUNNER_VERSION}/protoc-gen-php-grpc-${ROADRUNNER_VERSION}-linux-${arch}.tar.gz" \
    && tar -xzf "protoc-gen-php-grpc-${ROADRUNNER_VERSION}-linux-${arch}.tar.gz" \
    && mv "protoc-gen-php-grpc-${ROADRUNNER_VERSION}-linux-${arch}/protoc-gen-php-grpc" /usr/local/bin/protoc-gen-php-grpc \
    && chmod +x /usr/local/bin/protoc-gen-php-grpc \
    && rm "protoc-gen-php-grpc-${ROADRUNNER_VERSION}-linux-${arch}.tar.gz" \
    && apk del --no-cache wget tar

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
COPY --from=roadrunner /usr/bin/rr /usr/local/bin/rr

COPY ./.docker/chat-authz/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

FROM app_php_pre_composer AS app_php_composer
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"
COPY --from=composer_stage /usr/bin/composer /usr/local/bin/composer
RUN chown -R ${UID}:${GID} /usr/local/bin/composer
USER ${UID}

FROM app_php_pre_composer AS app_php_dev
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"
COPY --from=composer_stage /usr/bin/composer /usr/local/bin/composer
RUN chown -R ${UID}:${GID} /usr/local/bin/composer

RUN --mount=type=bind,from=php_ext_installer_stage,source=/usr/bin/install-php-extensions,target=/usr/local/bin/install-php-extensions \
    MAKEFLAGS="-j$(nproc)" install-php-extensions xdebug \
    && apk del --no-cache ${PHPIZE_DEPS} ${BUILD_DEPENDS}

COPY ./.docker/chat-authz/config/dev.ini ${PHP_INI_DIR}/conf.d/

USER ${UID}

FROM app_php_pre_composer AS app_php_test
ARG UID
ARG GID

ENV COMPOSER_HOME="/tmp/composer"
COPY --chown=${UID}:${GID} ./services/chat-authz .

COPY --from=composer_stage /usr/bin/composer /usr/local/bin/composer
RUN chown -R ${UID}:${GID} /usr/local/bin/composer
RUN composer install

COPY ./.docker/chat-authz/config/php.ini ${PHP_INI_DIR}/conf.d/

USER ${UID}
