//go:build !windows

package modapi

func ResolveRivalsPath() (string, error) {
	jeffm_only_supports_windows_please_use_windows()

	return "", nil
}
