#!/usr/bin/env bash
source "../utils.sh"

trap cleanup EXIT

section "Setting up"

cleanup
ck_env

export FOO=beginning

test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

cryptkeeper init "${TARGET_SHELL}" --standalone

ck_env

test_nempty "$CK_WATCH"
test_nempty "$CK_LAST"
test_nempty "$CK_REVERT"
test_eq "$FOO" "beginning"
test_eq "standalone" "$(cat .ckrc | jq -r .mode)"

section "Adding Secret from stdin"

echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

section "Overwriting secret"

echo "bar2" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar2"

section "Adding Secret from clipboard"

echo "blah!" | xclip -selection clipboard
cryptkeeper set CLIP -c
ck_env
test_eq "$CLIP" "blah!"

section "Moving deeper into the tree"

watch="${CK_WATCH}"
mkdir -p foo/bar/baz
pushd foo/bar/baz

ck_env
test_eq "$CLIP" "blah!"
test_eq "$CK_WATCH" "$watch"

popd
rm -rf ./foo

section "Moving outside of the tree"

pushd ..
ck_env

test_eq "$FOO" "beginning"
test_empty "$CLIP"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

popd
ck_env

section "Remove single value"

cryptkeeper remove FOO
ck_env
test_eq "$FOO" "beginning"

section "Removing multiple values"

echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

cryptkeeper remove FOO CLIP
ck_env
test_eq "$FOO" "beginning"
test_empty "$CLIP"

section "Removing config"

echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

cleanup
ck_env

test_eq "$FOO" "beginning"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"
unset FOO