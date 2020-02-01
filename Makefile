SHELL = /usr/bin/env bash

run_star:
	cd cmd/star && go install
	star new.star.py
