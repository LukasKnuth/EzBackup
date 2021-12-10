package restic

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Empty string = current directory.
const CURRENT_DIR = ""

func SafeInit() (err error) {
	err = execRestic([]string{"snapshots"}, CURRENT_DIR)
	var e *exec.ExitError
	if err != nil && errors.As(err, &e) {
		if e.ExitCode() != 0 {
			// No repo initialized yet!
			err = execRestic([]string{"init"}, CURRENT_DIR)
		}
	}
	// no error, the repo is initialized and accessible
	return
}

func Backup(pvcName string) error {
	workingDir := envFallback("BACKUP_TARGET_DIR", "/mnt/target")
	err := execRestic([]string{"backup", ".", "--no-cache", "--host", pvcName}, workingDir)
	return err
}

func Restore(pvcName string) error {
	workingDir := envFallback("RESTORE_TARGET_DIR", "/mnt/target")
	err := execRestic([]string{"restore", "latest", "--target", workingDir, "--no-cache", "--host", pvcName}, CURRENT_DIR)
	return err
}

func Version() error {
	return execRestic([]string{"version"}, CURRENT_DIR)
}

func execRestic(args []string, workingDir string) error {
	resticPath := envFallback("RESTIC_PATH", "/app/restic")
	command := exec.Cmd {
		Path: resticPath,
		Args: append([]string{resticPath}, args...),
		Env: nil, // Re-use environment from this process!
		Dir: workingDir,
	}
	var e *exec.ExitError
	output, err := command.CombinedOutput()
	fmt.Printf("--- Restic invocation ---\n")
	fmt.Printf("Path: %s\nArgs: %s\n", resticPath, args)
	if err != nil && errors.As(err, &e) {
		fmt.Printf("Exit Code: %d", e.ExitCode())
	} else {
		fmt.Println("Exit Code: 0")
	}
	fmt.Printf("--- Output Start ---\n%s\n--- Output End ---\n", output)
	return err
}

func envFallback(name string, fallback string) string {
	value, found := os.LookupEnv(name)
	if found {
		return value
	} else {
		return fallback
	}
}