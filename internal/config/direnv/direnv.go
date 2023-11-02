package direnv

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
)

const evalExportCmd = "eval $(cryptkeeper export %s)"
const envrc = ".envrc"

func EnvrcPath() string {
	envrcPath, err := fileutils.FindPathTo(envrc)
	if err != nil {
		logrus.Debugf("failed to find path to .envrc: %s", err)
		return ""
	}

	return envrcPath
}

func Reload() error {
	return exec.Command("direnv", "reload").Run()
}

func Integrate(envrcPath, shell string) error {
	cmd := evalStatement(shell)

	exists, err := fileutils.TextExistsInFile(envrcPath, cmd)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	envrcBytes, err := fileutils.ReadFile(envrcPath)
	if err != nil {
		return err
	}

	// Check if file ends in newline
	if envrcBytes[len(envrcBytes)-1] != '\n' {
		envrcBytes = append(envrcBytes, '\n')
	}

	// Append eval export to end of file
	envrcBytes = append(envrcBytes, []byte(cmd)...)

	return fileutils.WriteFile(envrcPath, envrcBytes, 0644)
}

func evalStatement(shell string) string {
	return fmt.Sprintf(evalExportCmd, shell)
}

func EvalStatement(shell string) string {
	return evalStatement(shell)
}
