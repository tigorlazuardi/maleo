#!/bin/bash
RAWTAG=$(git describe --abbrev=0 --tags)
CURTAG=${RAWTAG##*/}
CURTAG="${CURTAG/v/}"

# shellcheck disable=SC2162
IFS='.' read -a vers <<<"$CURTAG"

MAJ=${vers[0]}
MIN=${vers[1]}
PATCH=${vers[2]}
echo "== INFO: Current Tag: $MAJ.$MIN.$PATCH"

case "$DRONE_COMMIT_MESSAGE" in
*"#major"*)
	((MAJ += 1))
	MIN=0
	PATCH=0
	echo "== INFO: #major found, incrementing major version"
	;;
*"#minor"*)
	((MIN += 1))
	PATCH=0
	echo "== INFO: #minor found, incrementing minor version"
	;;
*)
	((PATCH += 1))
	echo "== INFO: incrementing Patch Version"
	;;
esac

NEWTAG="$MAJ.$MIN.$PATCH"

FILES=$(find . -name go.mod | grep -v '^\./go.mod$')

for f in $FILES; do
	PACKAGE=$(echo "$f" | cut -d/ -f2- | xargs dirname)
	PACKAGE_TAG="$PACKAGE/v$NEWTAG"
	echo "== INFO: Adding Tag: $PACKAGE_TAG"
	git tag "$PACKAGE_TAG"
done

echo "== INFO: Adding Tag: v$NEWTAG"
git tag v$NEWTAG
