ARG ALPINE_VERSION
ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS app_base
ARG GROUP_NAME=appuser
ARG USER_NAME=appuser
ARG UID
ARG GID
ARG PROTOC_VERSION

RUN apk update && apk add --no-cache "protoc=${PROTOC_VERSION}" "protobuf-dev"

RUN if getent group ${GID} >/dev/null 2>&1; then \
    EXISTING_GROUP=$(getent group ${GID} | cut -d: -f1); \
    adduser -u ${UID} --ingroup $EXISTING_GROUP -S -g "${USER_NAME}" ${USER_NAME}; \
    else \
    addgroup -g ${GID} -S ${GROUP_NAME}; \
    adduser -u ${UID} --ingroup ${GROUP_NAME} -S -g "${USER_NAME}" ${USER_NAME}; \
    fi \
    && mkdir -p /home/${USER_NAME} /app \
    && chown -R ${UID}:${GID} /home/${USER_NAME} /app

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

RUN chown -R ${UID}:${GID} /go

WORKDIR /app

COPY ./.docker/authn/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

FROM app_base AS app_go_dev
ARG UID
ARG GID

ENV APP_ENV=dev

USER ${UID}

FROM app_base AS app_go_test
ARG UID
ARG GID

ENV APP_ENV=test

COPY --chown=${UID}:${GID} ./services/authn .

USER ${UID}
