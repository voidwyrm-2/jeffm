package modapi

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func joinpath(parts ...string) string {
	return path.Clean(strings.Join(parts, string(os.PathSeparator)))
}

func pathbase(path string) string {
	s := strings.Split(path, string(os.PathSeparator))
	return s[len(s)-1]
}

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
		enabledMods, err := os.Create(path)
		enabledMods.Close()
		if err != nil {
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

func VerifyJeffFolder(home string) error {
	p := joinpath(home, jeffFolder)
	if dir, err := os.ReadDir(p); isFileNotFound(err, p) || (err == nil && len(dir) < 2) {
		return InitJeffFolder(home)
	} else if err != nil {
		return err
	}

	return nil
}

func InitJeffFolder(home string) error {
	jefff := joinpath(home, jeffFolder)

	ok, err := doesFileExist(jefff)
	if err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(jefff, os.ModeDir); err != nil {
			return err
		}
	}

	configf := joinpath(jefff, configFile)

	if ok {
		ok, err = doesFileExist(configf)
		if err != nil {
			return err
		}

		err = createDirIfNotExists(joinpath(jefff, "mods"))
		if err != nil {
			return err
		}

	}

	if !ok {
		conf, err := os.Create(configf)
		defer conf.Close()
		if err != nil {
			return err
		}

		rivalsPath, err := ResolveRivalsPath()
		if err != nil {
			return err
		}

		_, err = conf.WriteString(fmt.Sprintf(`rivalsPath = '%s'`, rivalsPath))
		if err != nil {
			return err
		}

		err = createDirIfNotExists(joinpath(rivalsPath, paksSubpath, "~mods"))
		if err != nil {
			return err
		}
	}

	modsf := joinpath(jefff, "mods")

	ok, err = doesFileExist(modsf)
	if err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(modsf, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}
