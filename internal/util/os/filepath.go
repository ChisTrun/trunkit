package os

import (
	"path/filepath"
	"strings"
)

func MergePath(prefix, suffix string) string {
	checkingSuffix := suffix
	for {
		if checkingSuffix == "." {
			break
		}

		newPrefix := strings.TrimSuffix(prefix, checkingSuffix)
		if newPrefix != prefix {
			return filepath.Clean(filepath.Join(newPrefix, suffix))
		}
		checkingSuffix = filepath.Dir(checkingSuffix)
	}

	return filepath.Clean(filepath.Join(prefix, suffix))
}

func TrimPathSuffix(path, suffix string) string {
	checkingSuffix := suffix
	for {
		if checkingSuffix == "." {
			break
		}

		newPrefix := strings.TrimSuffix(path, checkingSuffix)
		if newPrefix != path {
			return newPrefix
		}
		checkingSuffix = filepath.Dir(checkingSuffix)
	}

	return path
}
