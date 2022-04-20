#!/bin/bash

if [ -f cmd/semantic-release/semantic-release ]; then
    echo "removing semantic-release go build file"
    rm cmd/semantic-release/semantic-release;
fi