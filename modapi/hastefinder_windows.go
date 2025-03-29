//go:build windows

package modapi

import (
	"errors"
	"fmt"
	"syscall"
)

const (
	steamCommonPathA = `Program Files (x86)\Steam\steamapps\common`
	steamCommonPathB = `SteamLibrary\steamapps\common`
)

func ResolveRivalsPath() (string, error) {
	drives, err := getDrives()
	if err != nil {
		return "", err
	}

	fmt.Println(drives)

	for _, d := range drives {
		d += ":"

		pathA := joinpath(d, steamCommonPathA, rivalsFolderName)
		if ok, err := doesFileExist(pathA); err != nil {
			return "", err
		} else if ok {
			return pathA, nil
		}

		pathB := joinpath(d, steamCommonPathB, rivalsFolderName)
		if ok, err := doesFileExist(pathB); err != nil {
			return "", err
		} else if ok {
			return pathB, nil
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
