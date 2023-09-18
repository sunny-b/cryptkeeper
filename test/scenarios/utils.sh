has() {
  type -P "$1" &>/dev/null
}

ck_env() {
  eval "$(cryptkeeper env $TARGET_SHELL)"
}

ck_eval() {
  eval "$(cryptkeeper export $TARGET_SHELL)"
}

test_eq() {
  if [[ "$1" != "$2" ]]; then
    echo "FAILED: '$1' == '$2'"
    exit 1
  fi
}

test_neq() {
  if [[ "$1" == "$2" ]]; then
    echo "FAILED: '$1' != '$2'"
    exit 1
  fi
}

test_empty() {
  if [[ -n "$1" ]]; then
    echo "FAILED: '$1' is not empty"
    exit 1
  fi
}

test_nempty() {
  if [[ -z "$1" ]]; then
    echo "FAILED: '$1' is empty"
    exit 1
  fi
}

cleanup_ck() {
  rm -f .ck*
  rm -rf nest
}

cleanup_direnv() {
  rm -f .envrc
}

cleanup() {
  cleanup_ck
  cleanup_direnv
}

section() {
  echo
  echo "## $1 ##"
}

direnv_eval() {
  eval $(direnv export $TARGET_SHELL)
}