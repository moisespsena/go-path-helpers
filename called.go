package path_helpers

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetCalledFileNameSkip(skip int, abs ...bool) (pth string) {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		panic(errors.New("Information unavailable."))
	}
	if len(abs) == 0 || !abs[0] {
		for _, gp := range GOPATHS {
			if strings.HasPrefix(filename, gp) {
				return strings.TrimPrefix(
					strings.TrimPrefix(filename, gp)[1:],
					"src"+string(os.PathSeparator),
				)
			}
		}
		if filename[0] == os.PathSeparator {
			filename = filename[1:]
		}
		return TrimGoPathC(filename, "src")
	}
	return filename
}

func GetCalledFileName(abs ...bool) string {
	return GetCalledFileNameSkip(2, abs...)
}

func GetCalledDir(abs ...bool) string {
	file := GetCalledFileNameSkip(2, abs...)
	return filepath.Dir(file)
}

func GetCalledDirOrError(abs ...bool) string {
	file := GetCalledFileNameSkip(2, abs...)
	if file == "" {
		panic("Invalid dir.")
	}
	return filepath.Dir(file)
}
