BIN=go
OS=$(shell go env GOOS)

SRC=./cmd/instancer/instancer.go
DEST=./bin
DIST=./dist


all:
	mkdir -p $(DEST)
	[ "$(OS)" = "linux" ] && $(MAKE) linux-amd64 || true
	[ "$(OS)" = "darwin" ] && $(MAKE) darwin-arm64 || true

linux-amd64:
	GOOS=linux GOARCH=amd64 $(BIN) build -o $(DEST)/instancer-linux-amd64 $(SRC)

darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(BIN) build -o $(DEST)/instancer-darwin-arm64 $(SRC)


run:
	[ "$(OS)" = "linux" ] && $(MAKE) run-linux-amd64 || true
	[ "$(OS)" = "darwin" ] && $(MAKE) run-darwin-arm64 || true

run-linux-amd64: linux-amd64
	$(DEST)/instancer-linux-amd64

run-darwin-arm64: darwin-arm64
	$(DEST)/instancer-darwin-arm64


dist:
	mkdir -p $(DIST)
	[ "$(OS)" = "linux" ] && $(MAKE) dist-linux-amd64 || true
	[ "$(OS)" = "darwin" ] && $(MAKE) dist-darwin-arm64 || true
	cp atlas.hcl $(DIST)/atlas.hcl
	cp -r atlas $(DIST)/atlas
	cp entrypoint.sh $(DIST)/entrypoint.sh

dist-linux-amd64: linux-amd64
	cp $(DEST)/instancer-linux-amd64 $(DIST)/instancer

dist-darwin-arm64: darwin-arm64
	cp $(DEST)/instancer-darwin-arm64 $(DIST)/instancer


clean:
	rm -rf $(DEST)/*
	rm -rf $(DIST)/*
