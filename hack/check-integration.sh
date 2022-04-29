#!/bin/bash

set -o errexit
set -o nounset

go test -p 1 -tags integration -covermode=atomic -timeout 90s ./...
