#!/usr/bin/env fish

set -x TARGET_SHELL fish

cd (dirname (status filename))
set TEST_DIR $PWD
set -x PATH (dirname "$TEST_DIR"):(git rev-parse --show-toplevel)/bin:$PATH

# Reset the environment loading if any
set -e CK_REVERT
set -e CK_WATCH
set -e CK_LAST

function test_scenario
  cd "$TEST_DIR/scenarios/$argv[1]"
  set test_string "### Testing $argv[1] - SHELL: $TARGET_SHELL"
  if test -n "$TARGET_ENCRYPTION"
    set test_string "$test_string - ENCRYPTION: $TARGET_ENCRYPTION"
  end
  set test_string "$test_string ###"

  echo
  echo $test_string

  ./test.fish || exit 1

  echo "### Succeeded ###"

  cd /
  ck_env
end

### RUN ###

source "$TEST_DIR/scenarios/utils.fish"
ck_env

test_scenario base

set algorithms aes ecc serpent rsa
for algo in $algorithms
  set -x TARGET_ENCRYPTION $algo
  test_scenario encryption
end

test_scenario direnv
test_scenario direnv-integrate
test_scenario nested
test_scenario direnv-nested