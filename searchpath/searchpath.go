package searchpath

import (
	"github.com/kiyoptr/su/errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

// LookupInPaths will search the given paths for a file with given name and returns the full path to file/dir
// for example ~/.config/?/config.yml will replace ? with name and checks whether replaced path exists
func LookupInPaths(paths []string, name string) (found string, err error) {
	for _, p := range paths {
		fullPath := strings.ReplaceAll(p, "?", name)

		if strings.HasPrefix(fullPath, "~/") {
			fullPath, err = resolveToHomeDir(fullPath)
		} else if strings.HasPrefix(fullPath, "./") ||
			!strings.HasPrefix(fullPath, "/") {
			fullPath, err = resolveToCwd(strings.TrimPrefix(fullPath, "./"))
		}

		if err != nil {
			err = errors.Newif(err, "failed to resolve path %s", p)
			return
		}

		if _, err = os.Stat(fullPath); err == nil {
			found = fullPath
			return
		}
	}

	return
}

func resolveToHomeDir(path string) (result string, err error) {
	result, err = homedir.Expand(path)
	return
}

func resolveToCwd(path string) (result string, err error) {
	wd := ""
	wd, err = os.Getwd()
	if err != nil {
		return
	}

	result = filepath.Join(wd, path)
	return
}
