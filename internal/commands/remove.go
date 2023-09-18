package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/config/direnv"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
)

var Remove = &cobra.Command{
	Use:   "remove",
	Short: "Remove one or more key-value pairs from the config",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		envKeys := args
		if allKeys {
			envKeys = cfg.Env.Keys()
		}

		for _, key := range envKeys {
			delete(cfg.Env, key)
		}

		if cfg.Encryption.Type == crypt.ECC256 {
			keeper, err := crypt.NewKeeper(cfg.Encryption.Type, cfg.Encryption.KeyPath)
			if err != nil {
				logrus.
					WithError(err).
					Debug("failed to remove encryption keys")
			}

			if keeper != nil {
				for _, envKey := range envKeys {
					err = keeper.RemoveKey(envKey)
					if err != nil {
						logrus.
							WithError(err).
							Debugf("failed to remove encryption key for %s", envKey)
					}
				}
			}
		}

		err = config.Write(cfg)
		if err != nil {
			return err
		}

		if cfg.IsDirenvIntegrated() {
			err = direnv.Reload()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	Remove.Flags().BoolVarP(&allKeys, "all", "a", false, "Remove all key-value pairs")
}
