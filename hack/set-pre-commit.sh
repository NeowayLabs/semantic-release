#!/bin/bash
mkdir -p ~/.git-templates/hooks
cp -f ./hack//githooks/commit-msg ~/.git-templates/hooks
chmod +x ~/.git-templates/hooks/commit-msg
git init
