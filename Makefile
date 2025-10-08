#!/user/bin/make
SHELL = /bin/bash

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
	init-gw \
	init-user

init-backend-cu: \
	init-gw-cu \
	init-user-cu

restart-backend: \
	restart-gw \
	restart-user
up-backend: \
	up-gw \
	up-user
down-backend: \
	down-gw \
	down-user
down-clear-backend: \
	down-clear-gw \
	down-clear-user

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