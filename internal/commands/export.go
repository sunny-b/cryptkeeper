package commands

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"github.com/sunny-b/cryptkeeper/internal/shell"
)

var Export = &cobra.Command{
	Use:       "export",
	Short:     "Export decrypted environment variables",
	Hidden:    true,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}

		target := args[0]

		sh := shell.Detect(target)
		// .envrc doesn't support fish shell, so we use bash instead.
		if target == "fish" {
			sh = shell.Bash
		}

		keeper, err := crypt.NewKeeper(cfg.Encryption.Type, cfg.Encryption.KeyPath)
		if err != nil {
			return
		}

		exported := ""
		for key, cipher := range cfg.Env {
			decryptedValue, err := keeper.Decrypt(key, cipher)
			if err != nil {
				log.
					WithField("err", err.Error()).
					Warn("failed to decrypt value")
			}

			exported += sh.Export(key, decryptedValue)
		}

		exported += sh.Export(config.CKWatchEnvKey, cfg.Path)

		if _, ok := os.LookupEnv(config.CKWatchEnvKey); !ok {
			log.Info("cryptkeeper: loading")
		}

		fmt.Print(exported)
	},
}
