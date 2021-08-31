#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if ! command -v git &> /dev/null
then
    DEV_VER="development"
else
    DEV_VER="dev-$(git rev-parse --short HEAD)"
fi

VERSION=${VERSION:=$DEV_VER}

build() {
    EXT=""
    [[ $GOOS = "windows" ]] && EXT=".exe"
    echo "Building ${GOOS} ${GOARCH}"
    go build \
        -o ./bin/vipien-${GOOS}-${GOARCH}${EXT} \
        ./cmd/client
    
    if [[ $GOOS = "linux" ]]; then
        go build \
            -o ./bin/vipien-server-${GOOS}-${GOARCH}${EXT} \
            ./cmd/server
    fi
}

### test units
go test ./... -v -cover -race || exit 1

### multi arch binary build
GOOS=linux GOARCH=amd64 build
GOOS=darwin GOARCH=arm64 build

# GOOS=windows GOARCH=amd64 build