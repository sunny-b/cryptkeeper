function test_eq --argument-names a b
    if not test (count $argv) = 2
        echo "Error: " (count $argv) " arguments passed to `eq`: $argv"
        exit 1
    end

    if not test $a = $b
        printf "Error:\n - expected: %s\n -      got: %s\n" "$a" "$b"
        exit 1
    end
end

function test_neq --argument-names a b
    if not test (count $argv) = 2
        echo "Error: " (count $argv) " arguments passed to `neq`: $argv"
        exit 1
    end

    if test $a = $b
        printf "Error:\n - expected: %s\n -      got: %s\n" "$a" "$b"
        exit 1
    end
end

function has
  type -q $argv[1]
end

function ck_env
  eval (cryptkeeper env $TARGET_SHELL)
end

function ck_eval
  eval (cryptkeeper export $TARGET_SHELL)
end

function test_empty
  if test -n "$argv[1]"
    echo "FAILED: '$argv[1]' is not empty"
    exit 1
  end
end

function test_nempty
  if test -z "$argv[1]"
    echo "FAILED: '$argv[1]' is empty"
    exit 1
  end
end

function cleanup_ck
  set -l files .ck*
  if test -n "$files"
    rm -f $files
  end
end

function cleanup_direnv
  set -l file .envrc
  if test -n "$file"
    rm -f $file
  end
end

function cleanup
  cleanup_ck
  cleanup_direnv
end

function section
  echo
  echo "## $argv[1] ##"
end

function direnv_eval
  eval (direnv export $TARGET_SHELL)
end