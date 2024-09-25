#!/bin/bash

WASM_HEADLESS=off GOOS=js GOARCH=wasm go run .
