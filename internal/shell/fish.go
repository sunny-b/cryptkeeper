package shell

import (
	"fmt"
	"strings"
)

type fish struct{}

// Fish adds support for the fish shell as a host
var Fish Shell = fish{}

const fishHook = `
    function __cryptkeeper_export_eval --on-event fish_prompt;
        "{{.SelfPath}}" env fish | source;

        if test "$cryptkeeper_fish_mode" != "disable_arrow";
            function __cryptkeeper_cd_hook --on-variable PWD;
                if test "$cryptkeeper_fish_mode" = "eval_after_arrow";
                    set -g __cryptkeeper_export_again 0;
                else;
                    "{{.SelfPath}}" env fish | source;
                end;
            end;
        end;
    end;

    function __cryptkeeper_export_eval_2 --on-event fish_preexec;
        if set -q __cryptkeeper_export_again;
            set -e __cryptkeeper_export_again;
            "{{.SelfPath}}" env fish | source;
            echo;
        end;

        functions --erase __cryptkeeper_cd_hook;
    end;
`

func (sh fish) Shell() string {
	return "fish"
}

func (sh fish) RCFile() string {
	return "config.fish"
}

func (sh fish) Hook() string {
	return fishHook
}

func (sh fish) Export(key, value string) string {
	if key == "PATH" {
		command := "set -x -g PATH"
		for _, path := range strings.Split(value, ":") {
			command += " " + sh.escape(path)
		}
		return command + ";"
	}
	return "set -x -g " + sh.escape(key) + " " + sh.escape(value) + ";"
}

func (sh fish) Unset(key string) string {
	return "set -e -g " + sh.escape(key) + ";"
}
func (sh fish) ExportAll(e Export) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.Unset(key)
		} else {
			out += sh.Export(key, *value)
		}
	}
	return out
}

func (sh fish) escape(str string) string {
	in := []byte(str)
	out := "'"
	i := 0
	l := len(in)

	hex := func(char byte) {
		out += fmt.Sprintf("'\\X%02x'", char)
	}

	backslash := func(char byte) {
		out += string([]byte{BACKSLASH, char})
	}

	escaped := func(str string) {
		out += "'" + str + "'"
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch {
		case char == TAB:
			escaped(`\t`)
		case char == LF:
			escaped(`\n`)
		case char == CR:
			escaped(`\r`)
		case char <= US:
			hex(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char == BACKSLASH:
			backslash(char)
		case char <= TILDA:
			literal(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i++
	}

	out += "'"

	return out
}
