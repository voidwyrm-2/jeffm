package modapi

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "embed"
)

/*
Warning: this calls `os.Open`, so try to use `isFileNotFound` instead when possible
*/
func doesFileExist(path string) (bool, error) {
	f, err := os.Open(path)
	f.Close()
	if isFileNotFound(err, path) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func isFileNotFound(err error, path string) bool {
	if err == nil {
		return false
	}

	errStr := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(err.Error()), "Error:"))

	checkA := "open " + path + ": no such file or directory"
	checkB := "open " + path + ": The system cannot find the file specified."
	checkC := "open " + path + ": The system cannot find the path specified."

	return errStr == checkA || errStr == checkB || errStr == checkC
}

func getItemsInFolder(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	names := []string{}

	for _, file := range files {
		names = append(names, file.Name())
	}

	return names, nil
}

func createDirIfNotExists(path string) error {
	if ok, err := doesFileExist(path); err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(path, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}

func formatModName(name string, enabled bool) string {
	if enabled {
		return "[ENABLED] " + name
	}

	return "[DISABLED] " + name
}

func shellOut(command string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	if stderr.String() != "" {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func VerifyJeffFolder(home string) error {
	p := filepath.Join(home, jeffFolder)
	if dir, err := os.ReadDir(p); isFileNotFound(err, p) || (err == nil && len(dir) < 2) {
		return InitJeffFolder(home)
	} else if err != nil {
		return err
	}

	return nil
}

func InitJeffFolder(home string) error {
	jefff := filepath.Join(home, jeffFolder)

	ok, err := doesFileExist(jefff)
	if err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(jefff, os.ModeDir); err != nil {
			return err
		}
	}

	configf := filepath.Join(jefff, configFile)

	if ok {
		ok, err = doesFileExist(configf)
		if err != nil {
			return err
		}

		err = createDirIfNotExists(filepath.Join(jefff, "mods"))
		if err != nil {
			return err
		}

	}

	if !ok {
		rivalsPath, err := ResolveRivalsPath()
		if err != nil {
			return err
		}

		err = os.WriteFile(configf, []byte(fmt.Sprintf(`rivalsPath = '%s'`, rivalsPath)), 0o666)
		if err != nil {
			return err
		}

		err = createDirIfNotExists(filepath.Join(rivalsPath, paksSubpath, "~mods"))
		if err != nil {
			return err
		}
	}

	err = createDirIfNotExists(filepath.Join(jefff, "mods"))
	if err != nil {
		return err
	}

	profilef := filepath.Join(jefff, profileFolder)

	ok, err = doesFileExist(profilef)
	if err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(profilef, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}
