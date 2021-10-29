

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "main.version=$(GIT_VERSION)" -X "main.date=$(DATE)" -X "main.commit=$(GIT_COMMIT)" -X "main.builtBy=$(BUILDER)"

lint:
	go vet ./...
	revive -set_exit_status ./...


scan:
	gosec -no-fail -fmt sarif -out security.sarif ./...


build:
	go build -ldflags='$(VERSION_FLAGS)'