package shell

import "fmt"

type bash struct{}

// Bash shell instance
var Bash Shell = bash{}

const bashHook = `
_cryptkeeper_hook() {
  local previous_exit_status=$?;
  trap -- '' SIGINT;
  eval "$("{{.SelfPath}}" env bash)";
  trap - SIGINT;
  return $previous_exit_status;
};
if ! [[ "${PROMPT_COMMAND:-}" =~ _cryptkeeper_hook ]]; then
  PROMPT_COMMAND="_cryptkeeper_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
fi
`

func (sh bash) Shell() string {
	return "bash"
}

func (sh bash) RCFile() string {
	return ".bashrc"
}

func (sh bash) Hook() string {
	return bashHook
}

func (sh bash) Export(key, value string) string {
	return "export " + sh.escape(key) + "=" + sh.escape(value) + ";"
}

func (sh bash) ExportAll(e Export) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.Unset(key)
		} else {
			out += sh.Export(key, *value)
		}
	}
	return out
}

func (sh bash) Unset(key string) string {
	return "unset " + sh.escape(key) + ";"
}

func (sh bash) escape(str string) string {
	return BashEscape(str)
}

//nolint:stylecheck,revive
const (
	ACK           = 6
	TAB           = 9
	LF            = 10
	CR            = 13
	US            = 31
	SPACE         = 32
	AMPERSTAND    = 38
	SINGLE_QUOTE  = 39
	PLUS          = 43
	NINE          = 57
	QUESTION      = 63
	UPPERCASE_Z   = 90
	OPEN_BRACKET  = 91
	BACKSLASH     = 92
	UNDERSCORE    = 95
	CLOSE_BRACKET = 93
	BACKTICK      = 96
	LOWERCASE_Z   = 122
	TILDA         = 126
	DEL           = 127
)

// https://github.com/solidsnack/shell-escape/blob/master/Text/ShellEscape/Bash.hs
/*
A Bash escaped string. The strings are wrapped in @$\'...\'@ if any
bytes within them must be escaped; otherwise, they are left as is.
Newlines and other control characters are represented as ANSI escape
sequences. High bytes are represented as hex codes. Thus Bash escaped
strings will always fit on one line and never contain non-ASCII bytes.
*/
func BashEscape(str string) string {
	if str == "" {
		return "''"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)
	escape := false

	hex := func(char byte) {
		escape = true
		out += fmt.Sprintf("\\x%02x", char)
	}

	backslash := func(char byte) {
		escape = true
		out += string([]byte{BACKSLASH, char})
	}

	escaped := func(str string) {
		escape = true
		out += str
	}

	quoted := func(char byte) {
		escape = true
		out += string([]byte{char})
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch {
		case char == ACK:
			hex(char)
		case char == TAB:
			escaped(`\t`)
		case char == LF:
			escaped(`\n`)
		case char == CR:
			escaped(`\r`)
		case char <= US:
			hex(char)
		case char <= AMPERSTAND:
			quoted(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char <= PLUS:
			quoted(char)
		case char <= NINE:
			literal(char)
		case char <= QUESTION:
			quoted(char)
		case char <= UPPERCASE_Z:
			literal(char)
		case char == OPEN_BRACKET:
			quoted(char)
		case char == BACKSLASH:
			backslash(char)
		case char == UNDERSCORE:
			literal(char)
		case char <= CLOSE_BRACKET:
			quoted(char)
		case char <= BACKTICK:
			quoted(char)
		case char <= TILDA:
			quoted(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i++
	}

	if escape {
		out = "$'" + out + "'"
	}

	return out
}
