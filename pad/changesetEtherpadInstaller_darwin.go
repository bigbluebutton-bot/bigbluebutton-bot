// +build darwin

package pad

import (
	"fmt"
	"os"
	"os/exec"
)

func installEtherpad(folderPath string) error {
	// Source constants and useful functions.
	// Note: This might not have an effect in Go since environment variables sourced here won't be available to subsequent commands.
	cmd := exec.Command("sh", "-c", `. src/bin/functions.sh`)
	cmd.Dir = folderPath + "/etherpad-lite"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to source constants and functions: %v", err)
	}
	fmt.Println("sourced constants and functions")

	// Prepare the environment by installing dependencies.
	cmd = exec.Command("sh", "-c", `src/bin/installDeps.sh`)
	cmd.Dir = folderPath + "/etherpad-lite"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}
	fmt.Println("installed dependencies")

	return nil
}