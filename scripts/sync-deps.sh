#!/bin/bash

GOSUMDB=off go work sync
git add .

FILES=$(find . -name go.mod | grep -v '^\./go.mod$')

for f in $FILES; do
	PACKAGE=$(echo "$f" | cut -d/ -f2- | xargs dirname)
	echo "== INFO: running 'go mod tidy' for $PACKAGE"
	(cd "$PACKAGE" && GOSUMDB=off go mod tidy && (git add go.mod go.sum || true))
done
