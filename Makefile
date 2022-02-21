.PHONY: all check gocritic build test clean install uninstall

default: all

all: check test build

test:
	go test -v -coverprofile coverage.out ./...

check:
	pre-commit run --all-files

gocritic:
	gocritic check -enableAll ./...

build:
	# build Gtk UI client
	./scripts/build_gtk_app.sh

	# build the service
	CGO_ENABLED=0 go build -ldflags="-w -s" cmd/cpu-governor-svc/*.go

install:
	sudo ./scripts/install_service.sh

uninstall:
	sudo ./scripts/uninstall_service.sh

clean:
	@rm -fv ./cpu-indicator-gtk3
