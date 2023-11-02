package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/config/direnv"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
	"github.com/sunny-b/cryptkeeper/internal/shell"
)

var yesPrompt bool

var Direnv = &cobra.Command{
	Use:       "direnv",
	Short:     "Integrate with direnv",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		sh := shell.Detect(args[0])
		envrcPath := direnv.EnvrcPath()

		var integrate bool
		switch {
		case yesPrompt:
			if len(envrcPath) > 0 {
				integrate = true
			} else {
				logrus.Println("No .envrc file detected. Skipping direnv integration.")
			}
		case len(envrcPath) > 0:
			output := promptUserf(".envrc file detected at %s. Would you like to integrate with direnv? [Y/n]: ", envrcPath)
			if output == "y" || output == "yes" {
				integrate = true
			}
		}

		if integrate {
			exists, err := fileutils.TextExistsInFile(envrcPath, fmt.Sprintf("eval $(cryptkeeper export %s)", sh.Shell()))
			if err != nil || !exists {
				fmt.Println("Please add the following to your .envrc file:\n\n" + direnv.EvalStatement(sh.Shell()))
			}

			_ = direnv.Reload()

			cfg.Direnv = &config.Direnv{
				RCPath: envrcPath,
			}
		} else {
			return nil
		}

		err = config.Write(cfg)
		if err != nil {
			return fmt.Errorf("error writing config: %w", err)
		}

		return nil
	},
}

func init() {
	Direnv.Flags().BoolVarP(&yesPrompt, "yes", "y", false, "Automatic yes to prompts")
}
