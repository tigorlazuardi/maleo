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
	echo "== INFO: #major found, using next major version"
	;;
*"#minor"*)
	((MIN += 1))
	PATCH=0
	echo "== INFO: #minor found, using next minor version"
	;;
*)
	((PATCH += 1))
	echo "== INFO: using next Patch Version"
	;;
esac

NEWTAG="$MAJ.$MIN.$PATCH"

FILES=$(find . -name go.mod | grep -v '^\./go.mod$')

for f in $FILES; do
	echo "== INFO: Updating $PACKAGE to $NEWTAG"
	sed -i -r "s#(\s+github\.com/tigorlazuardi/maleo.*\sv).*#\1$NEWTAG#" "$f"
	git add "$f"
done

git commit -m "Bump Version to v$NEWTAG [CI SKIP]"
