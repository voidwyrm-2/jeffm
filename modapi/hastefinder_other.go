//go:build !windows

package modapi

func ResolveHastePath() (string, error) {
	jeffm_only_supports_windows_please_use_windows()

	return "", nil
}
