package commands

import (
	"bufio"
	"crypto/subtle"
	"errors"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"golang.org/x/term"
)

var Verify = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"check"},
	Short:   "Verifies the encrypted value of a specified secret key",
	Long:    "Verifies if the provided value matches the encrypted value of a specified secret key, without revealing the actual secret. Useful for confirming the integrity or correctness of a stored secret.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envKey, expectedValue := args[0], ""

		cfg, err := config.GetConfig()
		if err != nil {
			return errors.New("failed to get config")
		}

		cipher, ok := cfg.Env[envKey]
		if !ok {
			return fmt.Errorf("secret %s does not exist", envKey)
		}

		switch {
		case useClipboard:
			expectedValue, err = clipboard.ReadAll()
			if err != nil {
				logrus.WithError(err).Debug("failed to get value from clipboard")
			}
		case isInputPiped():
			reader := bufio.NewReader(os.Stdin)
			expectedValue, err = reader.ReadString('\n')
			if err != nil {
				logrus.WithError(err).Debug("failed to get piped in value")
			}

			expectedValue = expectedValue[:len(expectedValue)-1]
		}

		// Default to user passing in value if value isn't set.
		if expectedValue == "" {
			fmt.Print("Enter expected value (it won't be displayed): ")
			byteValue, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("error reading value: %w", err)
			}

			expectedValue = string(byteValue)
		}

		keeper, err := crypt.NewKeeper(cfg.Encryption.Type, cfg.Encryption.KeyPath)
		if err != nil {
			return err
		}

		decryptedValue, err := keeper.Decrypt(envKey, cipher)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		if subtle.ConstantTimeCompare([]byte(decryptedValue), []byte(expectedValue)) == 1 {
			fmt.Print("equal")
		} else {
			fmt.Print("not-equal")
		}

		return nil
	},
}

func init() {
	Verify.Flags().BoolVarP(&useClipboard, "clipboard", "c", false, "Read value from clipboard")
}
