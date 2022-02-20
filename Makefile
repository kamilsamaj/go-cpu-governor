.PHONY: all check gocritic build test clean

default: all

all: check test build

test:
	go test -v -coverprofile coverage.out ./...

check:
	pre-commit run --all-files

gocritic:
	gocritic check -enableAll ./...

build:
	./scripts/build.sh

clean:
	@rm -fv ./cpu-indicator-gtk3
