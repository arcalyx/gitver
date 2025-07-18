#!/bin/bash
set -e

VERSION=${VERSION:-"1.0.0"}
RELEASE_TYPE=${RELEASE_TYPE:-"stable"}

# Generate changelog
echo "Generating changelog..."
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
echo "Latest tag: ${LATEST_TAG:-none}"

if [[ "$RELEASE_TYPE" == "nightly" ]]; then
  echo "Generating changelog for nightly build (last 20 commits)"
  git log -n 20 --pretty=format:"* %s (%h)" > CHANGELOG.md
elif [[ "$RELEASE_TYPE" == "prerelease" ]]; then
  echo "Generating changelog for prerelease build"
  if [ -z "$LATEST_TAG" ]; then
    git log --pretty=format:"* %s (%h)" > CHANGELOG.md
  else
    git log ${LATEST_TAG}..HEAD --pretty=format:"* %s (%h)" > CHANGELOG.md
  fi
elif [ -z "$LATEST_TAG" ]; then
  echo "No previous tag found, including all commits in changelog"
  git log --pretty=format:"* %s (%h)" > CHANGELOG.md
else
  echo "Generating changelog since tag $LATEST_TAG"
  git log ${LATEST_TAG}..HEAD --pretty=format:"* %s (%h)" > CHANGELOG.md
fi

if [ ! -s CHANGELOG.md ]; then
  echo "No changes found, creating placeholder"
  echo "* No changes" > CHANGELOG.md
fi

echo "Changelog contents:"
cat CHANGELOG.md
