ARG POSTGRES_VERSION
ARG ALPINE_VERSION

FROM postgres:${POSTGRES_VERSION}-alpine${ALPINE_VERSION}

COPY ./.docker/postgres/scripts/init-databases.sh /docker-entrypoint-initdb.d/init-databases.sh
RUN chmod +x /docker-entrypoint-initdb.d/init-databases.sh
