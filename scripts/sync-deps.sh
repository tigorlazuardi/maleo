#!/bin/bash

echo "== INFO: Running 'go work sync'"
GOSUMDB=off go work sync
git add .

FILES=$(find . -name go.mod | grep -v '^\./go.mod$')

for f in $FILES; do
	PACKAGE=$(echo "$f" | cut -d/ -f2- | xargs dirname)
	echo "== INFO: running 'go mod tidy' for $PACKAGE"
	(cd "$PACKAGE" && go mod tidy && (git add go.mod go.sum 2>/dev/null || true))
done
