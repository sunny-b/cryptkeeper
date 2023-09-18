package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/config/direnv"
	"github.com/sunny-b/cryptkeeper/internal/crypt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	useClipboard bool
)

var Set = &cobra.Command{
	Use:     "set",
	Aliases: []string{"add"},
	Short:   "Set a new key-value pair",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], ""

		cfg, err := config.GetConfig()
		if err != nil {
			return errors.New("failed to get config")
		}

		switch {
		case useClipboard:
			value, err = clipboard.ReadAll()
			if err != nil {
				logrus.
					WithError(err).
					Warn("failed to get value from clipboard")
			}
		case isInputPiped():
			reader := bufio.NewReader(os.Stdin)
			value, err = reader.ReadString('\n')
			if err != nil {
				logrus.
					WithError(err).
					Warn("failed to get piped in value")
			}
		}

		// Default to user passing in value if value isn't set.
		if value == "" {
			fmt.Print("Enter value (it won't be displayed): ")
			byteValue, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("error reading value: %w", err)
			}

			value = string(byteValue)
		}

		value = strings.TrimSuffix(value, "\n")

		keeper, err := crypt.NewKeeper(cfg.Encryption.Type, cfg.Encryption.KeyPath)
		if err != nil {
			return err
		}

		// Encrypt the value
		encryptedValue, err := keeper.Encrypt(key, value)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		if cfg.Env == nil {
			cfg.Env = make(map[string]string)
		}

		// Add the key-value pair to the config
		cfg.Env[key] = encryptedValue

		err = config.Write(cfg)
		if err != nil {
			return errors.New("failed to write config")
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
	Set.Flags().BoolVarP(&useClipboard, "clipboard", "c", false, "Read value from clipboard")
}

func isInputPiped() bool {
	stat, err := os.Stdin.Stat()

	// If err is non-nil, we default to the user passing in the value.
	return err == nil && (stat.Mode()&os.ModeCharDevice) == 0
}
