#!/bin/bash
mkdir -p ~/.git-templates/hooks
cp -f ./hack//githooks/commit-msg ~/.git-templates/hooks
chmod +x ~/.git-templates/hooks/commit-msg
git init
echo "To activate commit message template in a pre-existing local repository, go to the project root folder and run `git init` command."
