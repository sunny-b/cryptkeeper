#!/usr/bin/env bash
source "../utils.sh"

trap cleanup EXIT

section "Setting up"

cleanup
ck_env

test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

cryptkeeper init "${TARGET_SHELL}" -e "${TARGET_ENCRYPTION}" --standalone

ck_env

test_nempty "$CK_WATCH"
test_nempty "$CK_LAST"
test_nempty "$CK_REVERT"

section "Adding Secret from stdin"

echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

section "Decrypting secret"

test_eq "$(cryptkeeper decrypt FOO)" "bar"

section "Verifying secret"

test_eq "$(echo 'bar' | cryptkeeper verify FOO)" "equal"
test_eq "$(echo 'false' | cryptkeeper verify FOO)" "not-equal"

section "Remove secret"

cryptkeeper remove FOO
ck_env
test_empty "$FOO"