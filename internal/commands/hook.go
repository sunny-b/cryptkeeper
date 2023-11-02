package commands

import (
	"html/template"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny-b/cryptkeeper/internal/shell"
)

type hookContext struct {
	SelfPath string
}

var Hook = &cobra.Command{
	Use:       "hook",
	Short:     "Prints the shell hook to stdout",
	Hidden:    true,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		selfPath, err := os.Executable()
		if err != nil {
			return err
		}

		// Convert Windows path if needed
		selfPath = strings.Replace(selfPath, "\\", "/", -1)

		ctx := hookContext{
			SelfPath: selfPath,
		}

		sh := shell.Detect(args[0])

		hookTemplate, err := template.New("hook").Parse(sh.Hook())
		if err != nil {
			return err
		}

		err = hookTemplate.Execute(os.Stdout, ctx)
		if err != nil {
			return err
		}

		return nil
	},
}
