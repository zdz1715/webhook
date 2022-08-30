package util

import (
	"os"
	"path/filepath"
)

func GetExecRootDir() string {
	execPath, _ := os.Executable()

	return filepath.Dir(execPath)
}
