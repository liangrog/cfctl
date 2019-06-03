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

fast:
	go build -o ${APPNAME} ${LDFLAGS}

linux:
	GOOS=linux GOARCH=386 go build -v ${LDFLAGS} -o ./target/linux_386/${APPNAME}
	GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/linux_amd64/${APPNAME}

darwin:
	GOOS=darwin GOARCH=386 go build -v ${LDFLAGS} -o ./target/darwin_386/${APPNAME}
	GOOS=darwin GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/darwin_amd64/${APPNAME}

windows:
	GOOS=windows GOARCH=386 go build -v ${LDFLAGS} -o ./target/windows_386/${APPNAME}.exe
	GOOS=windows GOARCH=amd64 go build -v ${LDFLAGS} -o ./target/windows_amd64/${APPNAME}.exe
