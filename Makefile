BINARY := xmpp-admin

.PHONY: all test vet build check clean

all: check build

test:
	go test ./...

vet:
	go vet ./...

build:
	CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o $(BINARY) .

check: test vet
	bash -n install.sh

clean:
	rm -f $(BINARY) coverage.out
