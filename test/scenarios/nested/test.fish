#!/usr/bin/env fish
source "../utils.fish"

section "Setting up"
cleanup
ck_env

test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

cryptkeeper init "$TARGET_SHELL"

ck_env

test_nempty "$CK_WATCH"
test_nempty "$CK_LAST"
test_nempty "$CK_REVERT"

section "Adding Secret from stdin"

echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

section "Nesting cryptkeeper config"

mkdir -p nest
cd nest
cleanup
ck_env

cryptkeeper init "$TARGET_SHELL"
ck_env
test_empty "$FOO"

section "Adding nested secret"

echo "baz" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "baz"

section "Moving up the tree"

cd ..
ck_env
test_eq "$FOO" "bar"

section "Moving back down the tree"

cd nest
ck_env
test_eq "$FOO" "baz"

section "Cleaning up nested"

cleanup
ck_env

test_eq "$FOO" "bar"

section "Cleaning up root"

cd ..
ck_env

test_eq "$FOO" "bar"

rm -rf nest
cleanup
ck_env

test_empty "$FOO"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"