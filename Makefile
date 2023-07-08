.PHONY: deps clean build

TAG=0.13.0
BINARY=inout
DIST_DIR=_dist
OS=darwin
ARCH=arm64
VERSION=tip
USER_BIN=${HOME}/bin

deps:
	go get -u ./...

clean: 
	rm -rf ${DIST_DIR}/*
	
build:
	./_bin/build.sh ${OS} ${VERSION} ${ARCH}

test:
	./_bin/test.sh

tag:
	./_bin/tag.sh ${TAG}

install:
	./_bin/install.sh ${OS} ${ARCH}

install_local: build
	chmod +x ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION}
	mv ${DIST_DIR}/${BINARY}_${OS}_${ARCH}_${VERSION} ${USER_BIN}/${BINARY}