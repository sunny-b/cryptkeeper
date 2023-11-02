package shell

// ZSH is a singleton instance of ZSH_T
type zsh struct{}

// Zsh adds support for the venerable Z shell.
var Zsh Shell = zsh{}

const zshHook = `
_cryptkeeper_hook() {
  trap -- '' SIGINT;
  eval "$("{{.SelfPath}}" env zsh)";
  trap - SIGINT;
}
typeset -ag precmd_functions;
if [[ -z "${precmd_functions[(r)_cryptkeeper_hook]+1}" ]]; then
  precmd_functions=( _cryptkeeper_hook ${precmd_functions[@]} )
fi
typeset -ag chpwd_functions;
if [[ -z "${chpwd_functions[(r)_cryptkeeper_hook]+1}" ]]; then
  chpwd_functions=( _cryptkeeper_hook ${chpwd_functions[@]} )
fi
`

func (sh zsh) Shell() string {
	return "zsh"
}

func (sh zsh) RCFile() string {
	return ".zshrc"
}

func (sh zsh) Hook() string {
	return zshHook
}

func (sh zsh) Export(key, value string) string {
	return "export " + sh.escape(key) + "=" + sh.escape(value) + ";"
}

func (sh zsh) Unset(key string) string {
	return "unset " + sh.escape(key) + ";"
}

func (sh zsh) escape(str string) string {
	return BashEscape(str)
}

func (sh zsh) ExportAll(e Export) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.Unset(key)
		} else {
			out += sh.Export(key, *value)
		}
	}
	return out
}
