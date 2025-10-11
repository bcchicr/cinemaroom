#!/user/bin/make
SHELL = /bin/bash

init: init-app
init-cu: init-app-cu
restart: restart-app
up: up-app
down: down-app
down-clear: down-clear-app

# App
init-app: \
	init-backend

init-app-cu: \
	init-backend-cu

restart-app: \
	restart-backend
up-app: \
	up-backend
down-app: \
	down-backend
down-clear-app: \
	down-clear-backend

# Backend
init-backend: \
	down-clear-backend \
    docker-pull-backend \
	docker-build-backend-current-user

init-backend-cu: \
	down-clear-backend \
    docker-pull-backend \
	docker-build-backend-current-user

restart-backend: \
	down-backend \
	up-backend
up-backend:
	docker-compose up -d gateway user authn chat playlist chat-authz postgres redis
down-backend:
	docker-compose down --remove-orphans gateway user authn chat playlist chat-authz postgres redis
down-clear-backend:
	docker-compose down --remove-orphans --volumes gateway user authn chat playlist chat-authz postgres redis

docker-pull-backend:
	docker-compose pull --ignore-pull-failures gateway user authn chat playlist chat-authz postgres redis
docker-build-backend:
	docker-compose build --pull gateway user authn chat playlist chat-authz postgres redis
docker-build-backend-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull gateway user authn chat playlist chat-authz postgres redis

# Gateway
init-gw: \
	down-clear-gw \
	docker-pull-gw \
	docker-build-gw

init-gw-cu:\
    down-clear-gw \
    docker-pull-gw \
	docker-build-gw-current-user

restart-gw: \
	up-gw \
	compose-du-gw \
	scan-cacheable-gw \
	down-gw \
	up-gw

up-gw:
	docker-compose up -d gateway
down-gw:
	docker-compose down --remove-orphans gateway
down-clear-gw:
	docker-compose down --remove-orphans --volumes gateway
compose-du-gw:
	docker-compose exec gateway composer du -o
scan-cacheable-gw:
	docker-compose exec gateway php ./bin/hyperf.php

protoc-gw:
	docker-compose exec gateway sh -c 'protoc \
		--plugin=protoc-gen-grpc=/usr/bin/grpc_php_plugin \
		--php_out=/app/generated \
		--grpc_out=/app/generated \
		--proto_path=/lib/proto \
		$$(find /lib/proto -name "*.proto")'
fmt-gw:
	docker-compose exec gateway php ./vendor/bin/php-cs-fixer fix

docker-pull-gw:
	docker-compose pull --ignore-pull-failures gateway
docker-build-gw:
	docker-compose build --pull gateway
docker-build-gw-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull gateway

# User
init-user: \
	down-clear-user \
	docker-pull-user \
	docker-build-user

init-user-cu:\
    down-clear-user \
    docker-pull-user \
	docker-build-user-current-user

restart-user: \
	down-user \
	up-user

up-user:
	docker-compose up -d user
down-user:
	docker-compose down --remove-orphans user
down-clear-user:
	docker-compose down --remove-orphans --volumes user

protoc-user:
	docker-compose exec user sh -c 'protoc \
		--plugin=protoc-gen-grpc=/usr/bin/protoc-gen-php-grpc \
		--php_out=/app/generated \
		--grpc_out=/app/generated \
		--proto_path=/lib/proto \
		$$(find /lib/proto -name "*.proto")'

docker-pull-user:
	docker-compose pull --ignore-pull-failures user
docker-build-user:
	docker-compose build --pull user
docker-build-user-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull user

# Authn
init-authn: \
	down-clear-authn \
	docker-pull-authn \
	docker-build-authn

init-authn-cu:\
    down-clear-authn \
    docker-pull-authn \
	docker-build-authn-current-user

restart-authn: \
	down-authn \
	up-authn

up-authn:
	docker-compose up -d authn
down-authn:
	docker-compose down --remove-orphans authn
down-clear-authn:
	docker-compose down --remove-orphans --volumes authn

protoc-authn:

docker-pull-authn:
	docker-compose pull --ignore-pull-failures authn
docker-build-authn:
	docker-compose build --pull authn
docker-build-authn-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull authn

# Chat
init-chat: \
	down-clear-chat \
	docker-pull-chat \
	docker-build-chat

init-chat-cu:\
    down-clear-chat \
    docker-pull-chat \
	docker-build-chat-current-user

restart-chat: \
	down-chat \
	up-chat

up-chat:
	docker-compose up -d chat
down-chat:
	docker-compose down --remove-orphans chat
down-clear-chat:
	docker-compose down --remove-orphans --volumes chat

protoc-chat:
	docker-compose exec chat sh -c 'protoc \
		--plugin=protoc-gen-grpc=/usr/bin/protoc-gen-php-grpc \
		--php_out=/app/generated \
		--grpc_out=/app/generated \
		--proto_path=/lib/proto \
		$$(find /lib/proto -name "*.proto")'

docker-pull-chat:
	docker-compose pull --ignore-pull-failures chat
docker-build-chat:
	docker-compose build --pull chat
docker-build-chat-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull chat

# Playlist
init-playlist: \
	down-clear-playlist \
	docker-pull-playlist \
	docker-build-playlist

init-playlist-cu:\
    down-clear-playlist \
    docker-pull-playlist \
	docker-build-playlist-current-user

restart-playlist: \
	down-playlist \
	up-playlist

up-playlist:
	docker-compose up -d playlist
down-playlist:
	docker-compose down --remove-orphans playlist
down-clear-playlist:
	docker-compose down --remove-orphans --volumes playlist

protoc-playlist:
	docker-compose exec playlist sh -c 'protoc \
		--plugin=protoc-gen-grpc=/usr/bin/protoc-gen-php-grpc \
		--php_out=/app/generated \
		--grpc_out=/app/generated \
		--proto_path=/lib/proto \
		$$(find /lib/proto -name "*.proto")'

docker-pull-playlist:
	docker-compose pull --ignore-pull-failures playlist
docker-build-playlist:
	docker-compose build --pull playlist
docker-build-playlist-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull playlist

# Chat-authz
init-chat-authz: \
	down-clear-chat-authz \
	docker-pull-chat-authz \
	docker-build-chat-authz

init-chat-authz-cu:\
    down-clear-chat-authz \
    docker-pull-chat-authz \
	docker-build-chat-authz-current-user

restart-chat-authz: \
	down-chat-authz \
	up-chat-authz

up-chat-authz:
	docker-compose up -d chat-authz
down-chat-authz:
	docker-compose down --remove-orphans chat-authz
down-clear-chat-authz:
	docker-compose down --remove-orphans --volumes chat-authz

protoc-chat-authz:
	docker-compose exec chat-authz sh -c 'protoc \
		--plugin=protoc-gen-grpc=/usr/bin/protoc-gen-php-grpc \
		--php_out=/app/generated \
		--grpc_out=/app/generated \
		--proto_path=/lib/proto \
		$$(find /lib/proto -name "*.proto")'

docker-pull-chat-authz:
	docker-compose pull --ignore-pull-failures chat-authz
docker-build-chat-authz:
	docker-compose build --pull chat-authz
docker-build-chat-authz-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull chat-authz

# Postgres
init-postgres: \
	down-clear-postgres \
	docker-pull-postgres \
	docker-build-postgres

init-postgres-cu:\
    down-clear-postgres \
    docker-pull-postgres \
	docker-build-postgres-current-user

restart-postgres: \
	down-postgres \
	up-postgres

up-postgres:
	docker-compose up -d postgres
down-postgres:
	docker-compose down --remove-orphans postgres
down-clear-postgres:
	docker-compose down --remove-orphans --volumes postgres

docker-pull-postgres:
	docker-compose pull --ignore-pull-failures postgres
docker-build-postgres:
	docker-compose build --pull postgres
docker-build-postgres-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull postgres

# Postgres
init-redis: \
	down-clear-redis \
	docker-pull-redis \
	docker-build-redis

init-redis-cu:\
    down-clear-redis \
    docker-pull-redis \
	docker-build-redis-current-user

restart-redis: \
	down-redis \
	up-redis

up-redis:
	docker-compose up -d redis
down-redis:
	docker-compose down --remove-orphans redis
down-clear-redis:
	docker-compose down --remove-orphans --volumes redis

docker-pull-redis:
	docker-compose pull --ignore-pull-failures redis
docker-build-redis:
	docker-compose build --pull redis
docker-build-redis-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull redis
