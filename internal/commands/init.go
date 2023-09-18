package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/config/direnv"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
	"github.com/sunny-b/cryptkeeper/internal/shell"
)

var (
	encryption string
	keyPath    string
	withDirenv bool
)

var Init = &cobra.Command{
	Use:       "init",
	Short:     "Initialize cryptkeeper",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath = fileutils.Clean(keyPath)
		configPath := fileutils.Clean(config.FileName())

		// Guard against overwriting key or config that already exist.
		if fileutils.FileExists(keyPath) {
			return fmt.Errorf("key file already exists at %s", keyPath)
		}
		if fileutils.FileExists(configPath) {
			return fmt.Errorf("config file already exists at %s", fileutils.Clean(config.FileName()))
		}

		var encType crypt.EncryptionType
		var err error
		switch strings.ToLower(encryption) {
		case "aes", "aes256", "aes-256":
			encType = crypt.AES256
		case "rsa", "rsa2048", "rsa-2048":
			encType = crypt.RSA2048
		case "ecc", "ecc256", "ecc-256":
			encType = crypt.ECC256
		case "serpent", "serpent256", "serpent-256":
			encType = crypt.Serpent256
		default:
			return crypt.ErrUnknownEncryptionType
		}

		err = crypt.GenerateKeys(encType, keyPath)
		if err != nil {
			return err
		}

		cfg := &config.Config{
			Encryption: config.Encryption{
				KeyPath: keyPath,
				Type:    encType,
			},
			Env:  make(config.Env),
			Path: configPath,
		}

		envrcPath := direnv.EnvrcPath()
		sh := shell.Detect(args[0])

		var integrate bool
		switch {
		case withDirenv && len(envrcPath) > 0:
			integrate = true
		case withDirenv:
			fmt.Println("No .envrc file detected. Skipping direnv integration.")
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

			cfg.Direnv = &config.Direnv{
				RCPath: envrcPath,
			}
		} else {
			rcPath, err := fileutils.FindPathTo(sh.RCFile())
			if err != nil || len(rcPath) == 0 {
				fmt.Printf("Please add the following to your %s file:\n\neval \"$(cryptkeeper hook %s)\"\n", sh.RCFile(), sh.Shell())
			}

			exists, err := fileutils.TextExistsInFile(rcPath, fmt.Sprintf(`eval "$(cryptkeeper hook %s)"`, sh.Shell()))
			if err != nil || !exists {
				fmt.Printf("Please add the following to your %s file:\n\neval \"$(cryptkeeper hook %s)\"\n", sh.RCFile(), sh.Shell())
			}
		}

		err = config.Init(cfg)
		if err != nil {
			return fmt.Errorf("error writing config: %w", err)
		}

		fmt.Printf("Initialized config in %s and key in %s\n", fileutils.Clean(config.FileName()), keyPath)

		return nil
	},
}

func init() {
	Init.Flags().StringVarP(&encryption, "encryption", "e", "aes256", "Type of encryption to use for encrypting/decrypting the secrets")
	Init.Flags().StringVarP(&keyPath, "key-path", "k", config.KeyFileName(), "File path to output generated encryption key. ")
	Init.Flags().BoolVarP(&withDirenv, "direnv", "d", false, "Integrate with direnv")
}

func promptUserf(prompt string, args ...any) string {
	fmt.Printf(prompt, args...)

	reader := bufio.NewReader(os.Stdin)
	output, _ := reader.ReadString('\n')

	return strings.TrimSpace(strings.ToLower(output))
}
