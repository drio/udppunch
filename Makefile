SERVER=punch-server
CLIENT=punch-client

PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

LDFLAGS=-ldflags '-s -w -extldflags "-static"' 

.PHONY: all build build_all clean vuln vuln/verbose test test/watch \
	coverage/html lint list copy

all: clean build_all list

copy: build_all
	rsync -avz -e ssh dist/punch-server-linux-amd64 atom:
	rsync -avz -e ssh dist/punch-client-linux-amd64 hs1:
	rsync -avz -e ssh dist/punch-client-linux-amd64 labs:

build:
	go build ${LDFLAGS} -o dist/${SERVER} server/server.go
	go build ${LDFLAGS} -o dist/${CLIENT} client/client.go

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o dist/$(SERVER)-$(GOOS)-$(GOARCH) server/server.go)))

	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o dist/$(CLIENT)-$(GOOS)-$(GOARCH) client/client.go)))

list:
	@ls -lachd dist/*

clean:
	@rm -rf dist c.out

vuln:
	govulncheck ./...

vuln/verbose:
	govulncheck -show verbose ./...

test:
	go test -v ./...

test/watch:
	@ls *.go server/*.go client/netx/*.go client/wg/*.go | \
		entr -c -s 'go test -failfast -v ./... && notify "💚" || notify "🛑"'

coverage/html:
	go test -v -cover -coverprofile=c.out
	go tool cover -html=c.out

lint:
	golangci-lint run
