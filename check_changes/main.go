package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path to scan>")
		return
	}
	pathToScan := os.Args[1]

	err := filepath.Walk(pathToScan, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if gitDirExists(path) {
				hasChanges, changes, err := gitHasChanges(path)
				if err != nil {
					return err
				}

				if hasChanges {
					color.Red("Repository with changes: %s", path)
					color.Yellow("%s", changes)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning directories: %v", err)
	}
}

func gitDirExists(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return !os.IsNotExist(err)
}

func gitHasChanges(path string) (bool, string, error) {
	cmd := exec.Command("git", "status", "-s")
	cmd.Dir = path
	output, err := cmd.Output()

	if err != nil {
		return false, "", err
	}

	return len(output) > 0, string(output), nil
}
