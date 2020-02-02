SHELL = /usr/bin/env bash


install:
	cd cmd/star && go install

run_star_requests: install
	cd examples && star requests.star.py

run_star_server: install
	cd examples && star server.star.py

run_star: install
	cd cmd/star && go install
	star
