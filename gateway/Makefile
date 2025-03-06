#!/usr/bin/make
SHELL = /bin/bash

init: \
	docker-down-clear \
	docker-pull \
	docker-build \
	docker-up

up: docker-up
down: docker-down
down-clear: docker-down-clear
restart: docker-down docker-up

utils:
	docker-compose exec -it utils sh
exec-phpunit:
	docker-compose exec utils ./vendor/bin/phpunit --colors=always
exec-psalm:
	docker-compose exec utils ./vendor/bin/psalm

set-up:
	cp -n ./.env.example ./.env && \
 	cp -n ./.app.env.example ./.app.env && \
 	cp -n ./.rr.example.yaml ./.rr.yaml
set-up-local:
	cp -n ./compose.local.example.yml ./compose.local.yml && \
	cp -n ./.docker/php/config/dev.example.ini ./.docker/php/config/dev.ini

docker-down:
	docker-compose down --remove-orphans
docker-down-clear:
	docker-compose down --remove-orphans --volumes

docker-pull:
	docker-compose pull --ignore-pull-failures

docker-build:
	docker-compose build --pull
docker-build-current-user:
	docker-compose build --build-arg UID=$$(id -u) --build-arg GID=$$(id -g) --pull

docker-up:
	docker-compose up -d
