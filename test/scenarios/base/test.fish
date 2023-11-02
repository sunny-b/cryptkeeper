#!/usr/bin/env fish
source ../utils.fish

function section
  echo
  echo "## $argv[1] ##"
end

# Setting up
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

# Adding Secret from stdin
section "Adding Secret from stdin"
echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"

# Overwriting secret
section "Overwriting secret"
echo "bar2" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar2"

# Adding Secret from clipboard
section "Adding Secret from clipboard"
echo "blah!" | xclip -selection clipboard
cryptkeeper set CLIP -c
ck_env
test_eq "$CLIP" "blah!"

# Moving deeper into the tree
section "Moving deeper into the tree"
set watch "$CK_WATCH"
mkdir -p foo/bar/baz
cd foo/bar/baz
ck_env
test_eq "$CLIP" "blah!"
test_eq "$CK_WATCH" "$watch"
cd -
rm -rf ./foo

# Moving outside of the tree
section "Moving outside of the tree"
cd ..
ck_env
test_empty "$FOO"
test_empty "$CLIP"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"
cd -
ck_env

# Remove single value
section "Remove single value"
cryptkeeper remove FOO
ck_env
test_empty "$FOO"

# Removing multiple values
section "Removing multiple values"
echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"
cryptkeeper remove FOO CLIP
ck_env
test_empty "$FOO"
test_empty "$CLIP"

# Removing config
section "Removing config"
echo "bar" | cryptkeeper set FOO
ck_env
test_eq "$FOO" "bar"
cleanup
ck_env
test_empty "$FOO"
test_empty "$CK_WATCH"
test_empty "$CK_LAST"
test_empty "$CK_REVERT"

cleanup