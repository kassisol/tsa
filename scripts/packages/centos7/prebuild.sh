#!/bin/bash

ROOTDIR=$(dirname $0)/../../..
cd $(dirname $0)

if [ -d "build" ]; then
	rm -rf build
fi
mkdir -p build/{bin,tsa,tsad}

cp ${ROOTDIR}/contrib/init/systemd/tsad.service build/
cp ${ROOTDIR}/bin/tsa build/bin/
cp ${ROOTDIR}/bin/tsad build/bin/

go run ${ROOTDIR}/gen/man/genman-tsa.go
cp -r /tmp/tsa/man build/tsa/

go run ${ROOTDIR}/gen/man/genman-tsad.go
cp -r /tmp/tsad/man build/tsad

go run ${ROOTDIR}/gen/shellcompletion/genshellcompletion-tsa.go
cp -r /tmp/tsa/shellcompletion build/tsa

go run ${ROOTDIR}/gen/shellcompletion/genshellcompletion-tsad.go
cp -r /tmp/tsad/shellcompletion build/tsad
