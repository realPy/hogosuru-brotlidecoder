#!/bin/bash

cd /go/src/brotlidec
tinygo build  --no-debug -target wasm  -o docs/brotlidecoder.wasm .
gzip -f -9 docs/brotlidecoder.wasm

