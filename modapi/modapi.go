package modapi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	unarr "github.com/gen2brain/go-unarr"

	"github.com/BurntSushi/toml"
)

//

const (
	modFileExtension = ".pak"
	rivalsFolderName = "MarvelRivals"
	jeffFolder       = "jeffmm"
	configFile       = "config.toml"
	profileFolder    = "profiles"
	paksSubpath      = `MarvelGame\Marvel\Content\Paks`

	lineEnding = "\r\n"
)

type ModHandler struct {
	rivalsPath, homePath, storePath, profilePath, loadedPath string
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

	_, err = toml.DecodeFile(filepath.Join(home, jeffFolder, configFile), &lmh)
	if err != nil {
		return ModHandler{}, err
	}

	mh := ModHandler{rivalsPath: lmh.RivalsPath, homePath: home, storePath: filepath.Join(home, jeffFolder, "mods"), profilePath: filepath.Join(home, jeffFolder, profileFolder), loadedPath: filepath.Join(lmh.RivalsPath, paksSubpath, "~mods")}
	if ok, err := doesFileExist(mh.loadedPath); err != nil {
		return ModHandler{}, err
	} else if !ok {
		pakPath := filepath.Join(lmh.RivalsPath, paksSubpath)
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
		if err := mh.InstallMod(mh.PathOfLoaded(mod), true, false); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) HomePath(subpath ...string) string {
	return filepath.Join(append([]string{mh.homePath}, subpath...)...)
}

func (mh ModHandler) Config() string {
	return mh.HomePath(jeffFolder, configFile)
}

func (mh ModHandler) PathOfStored(name string) string {
	return filepath.Join(mh.storePath, name)
}

func (mh ModHandler) PathOfLoaded(name string) string {
	return filepath.Join(mh.loadedPath, name)
}

func (mh ModHandler) GetRawMods() ([]string, error) {
	mods, err := getItemsInFolder(mh.storePath)
	if err != nil {
		return []string{}, err
	}

	return mods, nil
}

func (mh ModHandler) GetMods() ([]string, error) {
	mods, err := mh.GetRawMods()
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

func (mh ModHandler) GetProfiles() ([]string, error) {
	profiles, err := getItemsInFolder(mh.profilePath)
	if err != nil {
		return []string{}, err
	}

	return profiles, nil
}

/*
func (mh ModHandler) RepackMods(printProcess bool, modpaths ...string) error {
	for _, path := range modpaths {
		if err := mh.RepackMod(path, printProcess); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) RepackMod(modpath string, printProcess bool) error {
	repakPath := filepath.Join(pathUnbase(os.Args[0]), "repak.exe")

	ok, err := doesFileExist(repakPath)
	if err != nil {
		return err
	} else if ok {
		goto alreadyExists
	}

	err = os.WriteFile(repakPath, repakExecutable, 0o666)
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(repakPath)
	}()

alreadyExists:

	if printProcess {
		fmt.Println("unpacking...")
	}

	_, err = shellOut(repakPath, "unpack", modpath)
	if err != nil {
		return err
	}

	if printProcess {
		fmt.Println("repacking...")
	}

	if printProcess {
		fmt.Println("deleting old folder...")
	}

	_, err = shellOut(repakPath, "pack", strings.Split(modpath, ".")[0])
	if err != nil {
		return err
	}

	if printProcess {
		fmt.Println("finished process")
	}

	return os.RemoveAll(strings.Split(modpath, ".")[0])
}
*/

func (mh ModHandler) InstallMods(overwriteIfExists, printInstalls bool, modpaths ...string) error {
	for _, mp := range modpaths {
		if err := mh.InstallMod(mp, overwriteIfExists, printInstalls); err != nil {
			return err
		}
	}

	return nil
}

func (mh ModHandler) InstallRaw(name string, data []byte, overwriteIfExists bool) error {
	fl := os.O_RDWR | os.O_CREATE
	if overwriteIfExists {
		fl |= os.O_TRUNC
	}

	out, err := os.OpenFile(mh.PathOfStored(filepath.Base(name)), fl, 0o666)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = out.Write(data)
	return err
}

func (mh ModHandler) InstallReader(name string, r io.Reader, overwriteIfExists bool) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return mh.InstallRaw(name, data, overwriteIfExists)
}

func (mh ModHandler) InstallMod(modpath string, overwriteIfExists, printInstalls bool) error {
	modpath = path.Clean(modpath)
	ext := strings.ToLower(path.Ext(modpath))

	if slices.Contains([]string{".zip", ".rar", ".7z", ".tar"}, ext) {
		a, err := unarr.NewArchive(modpath)
		if err != nil {
			return err
		}

		defer a.Close()

		files, err := a.List()
		if err != nil {
			return err
		}

		for _, f := range files {
			if strings.ToLower(path.Ext(f)) == ".pak" {
				if printInstalls {
					fmt.Printf("found '%s'\n", f)
				}

				err := a.EntryFor(f)
				if err != nil {
					return err
				}

				data, err := a.ReadAll()
				if err != nil {
					return err
				}

				err = mh.InstallRaw(f, data, overwriteIfExists)
				if err != nil {
					return err
				}

				if printInstalls {
					fmt.Printf("installed '%s'\n", f)
				}
			}
		}
	} else if ext == ".pak" {
		pakr, err := os.Open(modpath)
		if err != nil {
			return err
		}

		defer pakr.Close()

		err = mh.InstallReader(filepath.Base(modpath), pakr, overwriteIfExists)
		if err != nil {
			return err
		}

		if printInstalls {
			fmt.Printf("installed '%s'\n", filepath.Base(modpath))
		}
	} else {
		return fmt.Errorf("'%s' is not an installable file", filepath.Base(modpath))
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
	name = filepath.Base(name)

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
