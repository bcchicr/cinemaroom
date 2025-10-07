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
	init-gw

init-backend-cu: \
	init-gw-cu

restart-backend: \
	restart-gw
up-backend: \
	up-gw
down-backend: \
	down-gw
down-clear-backend: \
	down-clear-gw

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
