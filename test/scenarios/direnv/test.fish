#!/usr/bin/env fish
source "../utils.fish"

section "Setting up"

cleanup_ck
ck_env

test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"
test_neq "$BAR" "baz"

cryptkeeper init "$TARGET_SHELL" -d

section "Testing direnv integration"

direnv allow .
direnv_eval
test_nempty "$CK_WATCH"
test_eq "$BAR" "baz"
test_nempty "$DIRENV_DIR"
test_nempty "$DIRENV_WATCHES"
test_nempty "$DIRENV_FILE"

section "Adding secret"
test_neq "$FOO" "bar"
echo "bar" | cryptkeeper set FOO
direnv_eval
test_eq "$FOO" "bar"

sleep 1

section "Updating secret"
test_neq "$FOO" "test"
echo "test" | cryptkeeper set FOO
direnv_eval
test_eq "$FOO" "test"

sleep 1

section "Leaving dir"
cd ..
direnv_eval
test_empty "$FOO"
cd -
direnv_eval
test_eq "$FOO" "test"

sleep 1

section "Removing secret"
cryptkeeper remove FOO
direnv_eval
test_neq "$FOO" "bar"

cleanup_ck