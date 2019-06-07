# Makefile

APPNAME=cfctl

VERSION_TAG=`git describe 2>/dev/null | cut -f 1 -d '-' 2>/dev/null`
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w \
    -X main.CommitHash=${COMMIT_HASH} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.Tag=${VERSION_TAG}"

all: clean test fast

test:
	go test -v ./...

clean:
	go clean
	rm -r target

fast:
	go build -o ${APPNAME} ${LDFLAGS}

linux:
	GOOS=linux GOARCH=386 go build -v ${LDFLAGS} -o ./target/${APPNAME}-linux-386
	GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/${APPNAME}-linux-amd64

darwin:
	GOOS=darwin GOARCH=386 go build -v ${LDFLAGS} -o ./target/${APPNAME}-darwin-386
	GOOS=darwin GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/${APPNAME}-darwin-amd64

windows:
	GOOS=windows GOARCH=386 go build -v ${LDFLAGS} -o ./target/${APPNAME}-windows-386.exe
	GOOS=windows GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/${APPNAME}-windows-amd64.exe

release: linux darwin windows
