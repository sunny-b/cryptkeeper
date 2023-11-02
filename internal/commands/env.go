package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"github.com/sunny-b/cryptkeeper/internal/envdiff"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
	"github.com/sunny-b/cryptkeeper/internal/shell"
	"github.com/sunny-b/cryptkeeper/internal/utils"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Env = &cobra.Command{
	Use:       "env",
	Short:     "Export or unset decrypted environment variables",
	Hidden:    true,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	Run: func(cmd *cobra.Command, args []string) {
		sh := shell.Detect(args[0])

		cfg, err := config.GetConfig()
		if err != nil {
			pathErr := new(*os.PathError)
			if (errors.Is(err, fileutils.ErrFileNotFound) || errors.As(err, pathErr)) && envVarExists(config.CKWatchEnvKey) && envVarExists(config.CKRevertEnvKey) {
				fmt.Print(unloadDiff(config.CKRevertEnvKey, sh))
			}
			return
		}

		if watchPath, ok := os.LookupEnv(config.CKWatchEnvKey); ok && watchPath != cfg.Path {
			diffString := unloadDiff(config.CKRevertEnvKey, sh)
			log.WithFields(log.Fields{
				"watch_path":  watchPath,
				"config_path": cfg.Path,
				"diff":        diffString,
			}).Debug("reverting env")

			fmt.Print(diffString)
		}

		if cfg.Direnv.Enabled() {
			// Don't run 'env' if direnv is enabled.
			fmt.Print(sh.Unset(config.CKRevertEnvKey) + sh.Unset(config.CKLastEnvKey))
			return
		}

		keeper, err := crypt.NewKeeper(cfg.Encryption.Type, cfg.Encryption.KeyPath)
		if err != nil {
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			return
		}

		// If the current working directory is a child of the directory containing the config file,
		// then we can export the environment variables. Otherwise, we'll just unset them.
		exportEnv, err := fileutils.IsChildDirOrSame(cwd, filepath.Dir(cfg.Path))
		if err != nil {
			return
		}

		if !exportEnv {
			fmt.Print(unloadDiff(config.CKRevertEnvKey, sh))
			return
		}

		lastEnv, err := envdiff.FetchLastEnv(config.CKLastEnvKey, keeper)
		if err != nil {
			log.WithError(err).Debug("failed to fetch last env")
			return
		}

		revertEnv, err := envdiff.FetchRevert(config.CKRevertEnvKey)
		if err != nil {
			log.WithError(err).Debug("failed to fetch revert enkv")
			return
		}

		currentEnv := cfg.Env
		err = currentEnv.Decrypt(keeper)
		if err != nil {
			log.WithError(err).Debug("failed to decrypt env")
			return
		}

		var firstLoad bool
		if loading() {
			log.Info("cryptkeeper: loading")
			firstLoad = true
		}

		log.WithFields(log.Fields{
			"CURRENT_ENV": currentEnv,
			"LAST_ENV":    lastEnv,
			"REVERT_ENV":  revertEnv,
		}).Debug("envs")

		if sameEnv(lastEnv, currentEnv) && len(revertEnv) == len(lastEnv) {
			log.Debug("cryptkeeper: no changes")
			if firstLoad {
				diffString := exportAllEnvs(cfg.Path, currentEnv, revertEnv, sh, keeper)
				log.WithField("diff", diffString).Debug("exporting")
				fmt.Print(diffString)
			}
			return
		}

		if out := diffStatus(envdiff.BuildEnvDiff(lastEnv, currentEnv)); out != "" {
			log.Infof("cryptkeeper: export %s", out)
		}

		newLast := lastEnv.Copy()
		for k, v := range currentEnv {
			if !utils.In(k, revertEnv) && !utils.In(k, lastEnv) {
				val, ok := os.LookupEnv(k)
				if ok {
					revertEnv[k] = utils.ToPtr(val)
				} else {
					revertEnv[k] = nil
				}
			}

			newLast[k] = v
		}

		for k := range lastEnv {
			if !utils.In(k, currentEnv) {
				delete(newLast, k)
			}
		}

		newRevertEnv := lo.MapValues(revertEnv, func(value *string, _ string) string {
			if value == nil {
				return ""
			}
			return *value
		})

		log.WithField("env", newRevertEnv).Debug("new revert env")

		diffString := envdiff.BuildEnvDiff(newRevertEnv, newLast).ToShell(sh)

		for k := range newRevertEnv {
			if !utils.In(k, newLast) {
				delete(revertEnv, k)
			}
		}

		diffString += exportAllEnvs(cfg.Path, currentEnv, revertEnv, sh, keeper)

		log.Debugf("env diff %s", diffString)
		fmt.Print(diffString)
	},
}

func envVarExists(envKey string) bool {
	_, ok := os.LookupEnv(envKey)
	return ok
}

func sameEnv(e1, e2 config.Env) bool {
	for k, v := range e1 {
		if e2[k] != v {
			return false
		}
	}

	for k, v := range e2 {
		if e1[k] != v {
			return false
		}
	}

	return true
}

//nolint:unparam
func unloadDiff(revertKey string, sh shell.Shell) string {
	revertEnv, err := envdiff.FetchRevert(revertKey)
	if err != nil {
		log.WithError(err).Debug("failed to fetch diff")
		return ""
	}
	if unloading() {
		log.Info("cryptkeeper: unloading")
	}

	var diffString string

	for k, v := range revertEnv {
		if v == nil {
			diffString += sh.Unset(k)
			os.Unsetenv(k)
		} else {
			diffString += sh.Export(k, *v)
			os.Setenv(k, *v)
		}
	}

	for _, envKey := range config.CKEnvKeys {
		diffString += sh.Unset(envKey)
		os.Unsetenv(envKey)
	}

	log.WithField("diff", diffString).Debug("unloading diff")

	return diffString
}

func exportAllEnvs(cfgPath string, currentEnv config.Env, revertEnv map[string]*string, sh shell.Shell, keeper *crypt.Keeper) string {
	diffString := ""
	encryptedDiff, err := keeper.Encrypt(config.CKLastEnvKey, config.Serialize(currentEnv))
	if err != nil {
		log.WithError(err).Debug("failed to encrypt diff")
	}

	if encryptedDiff != "" {
		diffString += sh.Export(config.CKLastEnvKey, encryptedDiff)
	}

	diffString += sh.Export(config.CKRevertEnvKey, config.Serialize(revertEnv))
	diffString += sh.Export(config.CKWatchEnvKey, cfgPath)

	return diffString
}

// Return a string of +/-/~ indicators of an environment diff
func diffStatus(oldDiff *envdiff.EnvDiff) string {
	if oldDiff.Any() {
		var out []string
		for key := range oldDiff.Revert {
			_, ok := oldDiff.Last[key]
			if !ok && !ckEnvKey(key) {
				out = append(out, "-"+key)
			}
		}

		for key := range oldDiff.Last {
			_, ok := oldDiff.Revert[key]
			if ckEnvKey(key) {
				continue
			}
			if ok {
				out = append(out, "~"+key)
			} else {
				out = append(out, "+"+key)
			}
		}

		sort.Strings(out)
		return strings.Join(out, " ")
	}

	return ""
}

func ckEnvKey(key string) bool {
	return strings.HasPrefix(key, "CK_")
}

func loading() bool {
	return !envVarExists(config.CKWatchEnvKey) && !envVarExists(config.CKRevertEnvKey) && !envVarExists(config.CKLastEnvKey)
}

func unloading() bool {
	return !loading()
}
