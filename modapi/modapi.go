package modapi

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	modFileExtension = ".pak"
	mrFolderName     = "MarvelRivals"
	jeffFolder       = "jeffmm"
	configFile       = "config.toml"
	enabledModsFile  = "enabled.txt"
	paksSubpath      = `MarvelGame\Marvel\Content\Paks`
)

type ModHandler struct {
	mrPath, homePath, loadedPath string
	enabledMods                  []string
}

func NewModHandler() (ModHandler, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ModHandler{}, err
	}

	lmh := struct {
		mrPath string
	}{}

	_, err = toml.DecodeFile(path.Join(home, jeffFolder, configFile), &lmh)
	if err != nil {
		return ModHandler{}, err
	}

	mh := ModHandler{mrPath: lmh.mrPath, enabledMods: []string{}, loadedPath: path.Join(lmh.mrPath, paksSubpath, "~mods")}
	if ok, err := doesFileExist(mh.loadedPath); err != nil {
		return ModHandler{}, err
	} else if !ok {
		pakPath := path.Join(lmh.mrPath, paksSubpath)
		if ok, err = doesFileExist(pakPath); err != nil {
			return ModHandler{}, err
		} else if !ok {
			return ModHandler{}, errors.New("Marvel Rivals is not installed, please install it to use Jeffm")
		}

		err = os.Mkdir(mh.loadedPath, os.ModeDir|0o666)
		if err != nil {
			return ModHandler{}, err
		}
	}

	return mh, nil
}

func (mh ModHandler) HomePath(subpath ...string) string {
	return path.Join(append([]string{mh.homePath}, subpath...)...)
}

func (mh ModHandler) VerifyRushFolder() error {
	p := mh.HomePath(jeffFolder)
	if dir, err := os.ReadDir(p); isFileNotFound(err, p) || (err == nil && len(dir) == 0) {
		return mh.InitRushFolder()
	} else if err != nil {
		return err
	}

	return nil
}

func (mh ModHandler) InitRushFolder() error {
	rushf := mh.HomePath(jeffFolder)

	ok, err := doesFileExist(rushf)
	if err != nil {
		return err
	} else if !ok {
		if err = os.Mkdir(rushf, os.ModeDir); err != nil {
			return err
		}
	}

	configf := path.Join(rushf, configFile)

	if ok {
		ok, err = doesFileExist(configf)
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

		mrPath, err := ResolveHastePath()
		if err != nil {
			return err
		}

		_, err = conf.WriteString(fmt.Sprintf(`mrPath = "%s"`, mrPath))
		if err != nil {
			return err
		}

		err = func() error {
			enabledMods, err := os.Create(enabledModsFile)
			enabledMods.Close()
			return err
		}()
		if err != nil {
			return err
		}
	}

	modsf := path.Join(rushf, "mods")

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

func (mh ModHandler) Config() string {
	return mh.HomePath(jeffFolder, configFile)
}

func (mh ModHandler) Close() error {
	return os.WriteFile(mh.HomePath(jeffFolder, "enabled.txt"), []byte(strings.Join(mh.enabledMods, "\n")), 0o666)
}

func (mh ModHandler) PathOfMod(name string) string {
	return path.Join(mh.homePath, jeffFolder, name)
}

func (mh ModHandler) PathOfEnabled(name string) string {
	return path.Join(mh.loadedPath, name)
}

func (mh ModHandler) GetMods() ([]string, error) {
	dir, err := os.ReadDir(path.Join(mh.homePath, jeffFolder))
	if err != nil {
		return []string{}, err
	}

	entries := []string{}

	for _, m := range dir {
		s := strings.Split(m.Name(), ".")
		entries = append(entries, formatModName(s[0], slices.Contains(mh.enabledMods, s[0])))
	}

	return entries, nil
}

func (mh ModHandler) InstallMods(modpaths ...string) error {
	for _, mp := range modpaths {
		if err := mh.InstallMod(mp); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) InstallMod(modpath string) error {
	ext := path.Ext(modpath)
	if ext == ".zip" {
		z, err := zip.OpenReader(modpath)
		defer z.Close()
		if err != nil {
			return err
		}

		for _, f := range z.File {
			ext = path.Ext(f.Name)
			if ext == ".pak" {
				fr, err := f.OpenRaw()
				if err != nil {
					return err
				}

				content, err := io.ReadAll(fr)
				if err != nil {
					return err
				}

				err = func() error {
					out, err := os.Create(mh.PathOfEnabled(path.Base(f.Name)))
					defer out.Close()
					if err != nil {
						return err
					}

					_, err = out.Write(content)
					return err
				}()
				if err != nil {
					return err
				}
			}
		}
	} else if ext == ".pak" {
		content, err := os.ReadFile(modpath)
		if err != nil {
			return err
		}

		out, err := os.Create(path.Join(mh.loadedPath, path.Base(modpath)))
		defer out.Close()
		if err != nil {
			return err
		}

		_, err = out.Write(content)
		return err
	}

	return fmt.Errorf("'%s' is not an installable file", path.Base(modpath))
}

func (mh *ModHandler) EnableMods(names ...string) error {
	for _, name := range names {
		if err := mh.EnableMod(name); err != nil {
			return err
		}
	}

	return nil
}

func (mh *ModHandler) EnableMod(name string) error {
	pakPath := mh.PathOfMod(name)

	pakContent, err := os.ReadFile(pakPath)
	if err != nil {
		return err
	}

	return os.WriteFile(mh.PathOfEnabled(name), pakContent, 0o666)
}

func (mh *ModHandler) DisableMods(names ...string) error {
	for _, name := range names {
		if err := mh.DisableMod(name); err != nil {
			return err
		}
	}

	return nil
}

func (mh *ModHandler) DisableMod(name string) error {
	i := slices.Index(mh.enabledMods, name)
	if i == -1 {
		return nil
	}

	mh.enabledMods = slices.Delete(mh.enabledMods, i, i+1)

	return os.Remove(mh.PathOfEnabled(name))
}
