#!/bin/bash

# Copyright (2021) Cobalt Speech and Language Inc.

rm -rf opus 
git clone https://github.com/xiph/opus.git
mkdir -p opus/build && cd opus/build
echo CMAKE_TOOLCHAIN_FILE $CMAKE_TOOLCHAIN_FILE
if [[ "${GOARCH}" = "arm" ]]; then
    CMAKE_C_FLAGS="-mfpu=neon"
fi
cmake -DCMAKE_TOOLCHAIN_FILE=${CMAKE_TOOLCHAIN_FILE} -DCMAKE_C_FLAGS="${CMAKE_C_FLAGS}" ..
make 



