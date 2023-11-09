#!/usr/bin/env fish
source "../utils.fish"

function reeval
    direnv_eval
end

section "Setting up"

cleanup
ck_env
reeval

echo "eval \$(cryptkeeper export $TARGET_SHELL)" > .envrc

sleep 1
reeval
direnv allow .

test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

cryptkeeper init "$TARGET_SHELL"
reeval

section "Adding Secret from stdin"

echo "bar" | cryptkeeper set FOO
reeval
test_eq "$FOO" "bar"

section "Nesting direnv"

mkdir -p nest
cd nest
cleanup
reeval

echo "export BAR=baz" > .envrc
echo "eval \$(cryptkeeper export $TARGET_SHELL)" >> .envrc

direnv allow .
reeval
test_nempty "$CK_WATCH"
test_eq "$BAR" "baz"
test_nempty "$DIRENV_DIR"
test_nempty "$DIRENV_WATCHES"
test_nempty "$DIRENV_FILE"

section "Nesting cryptkeeper config"

cryptkeeper init "$TARGET_SHELL"
reeval
test_eq "$FOO" "bar"

sleep 1

section "Adding nested secret"

echo "baz" | cryptkeeper set FOO
reeval
test_eq "$FOO" "baz"

sleep 1

section "Moving up the tree"

cd ..
reeval
test_eq "$FOO" "bar"

sleep 1

section "Moving back down the tree"

cd nest
reeval
test_eq "$FOO" "baz"
test_eq "$BAR" "baz"

sleep 1

section "Cleaning up nested"

cleanup
reeval
test_eq "$FOO" "bar"

sleep 1

section "Cleaning up root"

cd ..
reeval
test_eq "$FOO" "bar"

cleanup
reeval

section "Testing cleanup"

test_empty "$FOO"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"
test_empty "$DIRENV_DIR"
test_empty "$DIRENV_WATCHES"
test_empty "$DIRENV_FILE"
