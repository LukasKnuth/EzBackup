package operations

import (
	"fmt"
	"errors"

	"github.com/LukasKnuth/EzBackup/restic"
)

func Backup(pvcName string) error {
	err := restic.SafeInit()
	if err != nil {
		return errors.New("Restic could not initialize repository. See output for details...")
	}
	err = restic.Backup(pvcName)
	if err != nil {
		return errors.New("Restic could not start/finish backup. See output for details...")
	}
	fmt.Println("TODO: run RETAIN via restic")
	fmt.Println("TODO: run PURGE via restic")
	return nil
}

func Restore(pvcName string) error {
	err := restic.Restore(pvcName)
	if err != nil {
		return errors.New("Restic could not start/complete restore. See output for details...")
	}
	return nil
}