#!/usr/bin/env bash

mkdir -p output
cd output
go build ..

if [ -x ./apollo-client ]; then
    echo "Build success."
else
    echo "Build failed."
fi
