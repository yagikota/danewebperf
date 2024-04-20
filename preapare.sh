#!/bin/bash

prepare_docker_images() {
	unbound_dir="pageloadtime/docker/unbound"
	unbound_cache_image="unbound:with-cache"
	unbound_no_cache_image="unbound:without-cache"
	firefox_har_letsdane_dir="pageloadtime/docker/firefox/letsdane"
	letsdane_image="letsdane"
	firefox_har_dir="pageloadtime/docker/firefox"
	firefox_har_image="firefox"

	echo "Building docker images..."

	echo "Building unbound docker image..."
	cd "$unbound_dir"
	docker build -t "$unbound_cache_image" -f Dockerfile.unbound-with-cache .
	docker build -t "$unbound_no_cache_image" -f Dockerfile.unbound-without-cache .
	cd "$OLDPWD"

	echo "Building firefox docker image..."
	cd "$firefox_har_letsdane_dir"
	docker build -t "$letsdane_image" -f Dockerfile.letsdane .
	cd "$OLDPWD"

	echo "Building firefox docker image..."
	cd "$firefox_har_dir"
	docker build -t "$firefox_har_image" -f Dockerfile.firefox .
	cd "$OLDPWD"
}

prepare_docker_images
