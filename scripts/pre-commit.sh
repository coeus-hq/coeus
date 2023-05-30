#!/usr/bin/env bash

ROOT=$(git rev-parse --show-toplevel)
echo "Running pre-commit hook"

branch="$(git rev-parse --abbrev-ref HEAD)"
if [ "$branch" = "master" ]; then
  echo "master branch commit is blocked"
  exit 1
fi

go test ./...

# $? stores exit value of the last command
if [ $? -ne 0 ]; then
 echo "Tests must pass before commit!"
 exit 1
fi
