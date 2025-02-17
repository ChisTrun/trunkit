package os

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Exec(commands []string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", strings.Join(commands, ";"))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	return out.String(), nil
}

func CreateDirIfNotExists(dirPath string, perm os.FileMode) error {
	if _, err := os.Stat(dirPath); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("stat dir failed", dirPath, err)
			return err
		}

		if err = os.MkdirAll(dirPath, perm); err != nil {
			fmt.Println("make dir failed", dirPath, err)
			return err
		}
	}

	return nil
}
