#!/bin/bash

ROOTDIR=$(dirname $0)/../../..
cd $(dirname $0)

if [ -d "build" ]; then
	rm -rf build
fi
mkdir -p build

cp ${ROOTDIR}/contrib/init/systemd/tsa.service build/
cp ${ROOTDIR}/bin/tsa build/

go run ${ROOTDIR}/gen/man/genman.go
cp -r /tmp/tsa/man build/

go run ${ROOTDIR}/gen/shellcompletion/genshellcompletion.go
cp -r /tmp/tsa/shellcompletion build/
