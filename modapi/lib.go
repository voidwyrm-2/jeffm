package modapi

import (
	"os"
	"strings"
)

/*
Warning: this calls `os.Open`, so try to use `isFileNotFound` instead when possible
*/
func doesFileExist(path string) (bool, error) {
	f, err := os.Open(path)
	defer f.Close()
	if isFileNotFound(err, path) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func isFileNotFound(err error, file string) bool {
	if err == nil {
		return false
	}

	s := strings.Split(err.Error(), ":")
	if len(s) < 3 {
		return false
	}

	if err.Error() == "Error: open "+file+": no such file or directory" {
		return true
	} else {
		return strings.HasPrefix(s[1], "open ") && strings.HasSuffix(s[2], file+": no such file or directory")
	}
}

func formatModName(name string, enabled bool) string {
	if enabled {
		return "[ENABLED] " + name
	}

	return "[DISABLED] " + name
}
