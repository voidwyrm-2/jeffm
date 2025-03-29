package modapi

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"

	"github.com/BurntSushi/toml"
)

const (
	modFileExtension = ".pak"
	rivalsFolderName = "MarvelRivals"
	jeffFolder       = "jeffmm"
	configFile       = "config.toml"
	paksSubpath      = `MarvelGame\Marvel\Content\Paks`
)

type ModHandler struct {
	rivalsPath, homePath, storePath, loadedPath string
}

func NewModHandler() (ModHandler, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ModHandler{}, err
	}

	err = VerifyJeffFolder(home)
	if err != nil {
		return ModHandler{}, err
	}

	lmh := struct {
		RivalsPath string
	}{}

	_, err = toml.DecodeFile(joinpath(home, jeffFolder, configFile), &lmh)
	if err != nil {
		return ModHandler{}, err
	}

	mh := ModHandler{rivalsPath: lmh.RivalsPath, homePath: home, storePath: joinpath(home, jeffFolder, "mods"), loadedPath: joinpath(lmh.RivalsPath, paksSubpath, "~mods")}
	if ok, err := doesFileExist(mh.loadedPath); err != nil {
		return ModHandler{}, err
	} else if !ok {
		pakPath := joinpath(lmh.RivalsPath, paksSubpath)
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

	return mh, mh.storeAlreadyLoaded()
}

func (mh ModHandler) storeAlreadyLoaded() error {
	mods, err := getItemsInFolder(mh.loadedPath)
	if err != nil {
		return err
	}

	for _, mod := range mods {
		if err := mh.InstallMod(mh.PathOfLoaded(mod), true); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) HomePath(subpath ...string) string {
	return joinpath(append([]string{mh.homePath}, subpath...)...)
}

func (mh ModHandler) Config() string {
	return mh.HomePath(jeffFolder, configFile)
}

func (mh ModHandler) PathOfStored(name string) string {
	return joinpath(mh.storePath, name)
}

func (mh ModHandler) PathOfLoaded(name string) string {
	return joinpath(mh.loadedPath, name)
}

func (mh ModHandler) GetMods() ([]string, error) {
	mods, err := getItemsInFolder(mh.storePath)
	if err != nil {
		return []string{}, err
	}

	enabledMods, err := getItemsInFolder(mh.loadedPath)
	if err != nil {
		return []string{}, err
	}

	entries := []string{}

	for _, mod := range mods {
		entries = append(entries, formatModName(mod, slices.Contains(enabledMods, mod)))
	}

	return entries, nil
}

func (mh ModHandler) InstallMods(ignoreIfExists bool, modpaths ...string) error {
	for _, mp := range modpaths {
		if err := mh.InstallMod(mp, ignoreIfExists); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) InstallReader(name string, r io.Reader, ignoreIfExists bool) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	fl := os.O_RDWR | os.O_CREATE
	if !ignoreIfExists {
		fl |= os.O_TRUNC
	}

	out, err := os.OpenFile(mh.PathOfStored(pathbase(name)), fl, 0o666)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = out.Write(content)
	return err
}

func (mh ModHandler) InstallMod(modpath string, ignoreIfExists bool) error {
	modpath = path.Clean(modpath)
	ext := path.Ext(modpath)

	if ext == ".zip" {
		z, err := zip.OpenReader(modpath)
		defer z.Close()
		if err != nil {
			return err
		}

		for _, f := range z.File {
			ext = path.Ext(f.Name)
			if f.FileInfo().IsDir() {
			} else if ext == ".pak" {
				fr, err := f.OpenRaw()
				if err != nil {
					return err
				}

				err = mh.InstallReader(f.Name, fr, ignoreIfExists)
				if err != nil {
					return err
				}
			}
		}
	} else if ext == ".pak" {
		pakr, err := os.Open(modpath)
		defer pakr.Close()
		if err != nil {
			return err
		}

		return mh.InstallReader(pathbase(modpath), pakr, ignoreIfExists)
	} else {
		return fmt.Errorf("'%s' is not an installable file", pathbase(modpath))
	}

	return nil
}

func (mh ModHandler) UninstallMods(names ...string) error {
	for _, name := range names {
		if err := mh.UninstallMod(name); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) UninstallMod(name string) error {
	name = pathbase(name)

	storedMods, err := getItemsInFolder(mh.storePath)
	if err != nil {
		return err
	}

	loadedMods, err := getItemsInFolder(mh.loadedPath)
	if err != nil {
		return err
	}

	isStored := slices.Contains(storedMods, name)
	if isStored {
		if err = os.Remove(mh.PathOfStored(name)); err != nil {
			return err
		}
	}

	isLoaded := slices.Contains(loadedMods, name)
	if isLoaded {
		if err = os.Remove(mh.PathOfLoaded(name)); err != nil {
			return err
		}
	}

	if !isStored && !isLoaded {
		return fmt.Errorf("the mod '%s' does not exist", name)
	}

	return nil
}

func (mh ModHandler) EnableMods(names ...string) error {
	for _, name := range names {
		if err := mh.EnableMod(name); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) EnableMod(name string) error {
	pakPath := mh.PathOfStored(name)

	pakContent, err := os.ReadFile(pakPath)
	if err != nil {
		return err
	}

	return os.WriteFile(mh.PathOfLoaded(name), pakContent, 0o666)
}

func (mh ModHandler) DisableMods(names ...string) error {
	for _, name := range names {
		if err := mh.DisableMod(name); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) DisableMod(name string) error {
	return os.Remove(mh.PathOfLoaded(name))
}
