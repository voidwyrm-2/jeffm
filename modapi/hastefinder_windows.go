//go:build windows

package modapi

import (
	"errors"
	"fmt"
	"path"
	"syscall"
)

const (
	steamCommonPathA = `Program Files (x86)\Steam\steamapps\common`
	steamCommonPathB = `SteamLibrary\steamapps\common`
)

func ResolveHastePath() (string, error) {
	drives, err := getDrives()
	if err != nil {
		return "", err
	}

	for _, d := range drives {
		p := path.Join(d+":", steamCommonPathA, mrFolderName)
		if ok, err := doesFileExist(p); err != nil {
			return "", err
		} else if ok {
			return p, nil
		}
	}

	return "", errors.New("unable to find the Marvel Rivals folder")
}

func getDrives() ([]string, error) {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

	if ret, _, callErr := syscall.SyscallN(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); callErr != 0 {
		return []string{}, fmt.Errorf("error code %d while getting drives", callErr)
	} else {
		return bitsToDrives(uint32(ret)), nil
	}
}

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}
