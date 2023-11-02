package commands

import (
	"fmt"

	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/crypt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var allKeys bool
var withKey bool

var Decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt key(s) and print to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.GetConfig()
		if err != nil {
			return err
		}

		keeper, err := crypt.NewKeeper(config.Encryption.Type, config.Encryption.KeyPath)
		if err != nil {
			return err
		}

		var envKeys []string
		switch allKeys {
		case true:
			envKeys = config.Env.Keys()
		default:
			err = validateEnv(config.Env, args)
			if err != nil {
				return err
			}

			envKeys = args
		}

		for _, key := range envKeys {
			value, err := keeper.Decrypt(key, config.Env[key])
			if err != nil {
				log.
					WithField("err", err.Error()).
					Warn("failed to decrypt value")
			}

			if withKey {
				fmt.Printf("%s=%s\n", key, value)
			} else {
				fmt.Printf("%s\n", value)
			}
		}

		return nil
	},
}

func init() {
	Decrypt.Flags().BoolVarP(&allKeys, "all", "a", false, "Print all decrypted values")
	Decrypt.Flags().BoolVarP(&withKey, "with-key", "k", false, "Print out the key with the value in the format: KEY=VALUE")
}

func validateEnv(env map[string]string, keys []string) error {
	for _, key := range keys {
		if _, ok := env[key]; !ok {
			return fmt.Errorf("key %s not found", key)
		}
	}

	return nil
}
