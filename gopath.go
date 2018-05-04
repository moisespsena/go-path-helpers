package path_helpers

import (
	"os"
	"path"
	"strings"
	"go/build"
	"github.com/mitchellh/go-homedir"
)

var GOPATH string
var GOPATHS []string

func init() {
	var (
		err error
		pth string
		ok  bool
	)

	paths := make(map[string]int)

	for _, pth = range strings.Split(os.Getenv("GOPATH"), ":") {
		if pth != "" {
			if _, ok = paths[pth]; !ok {
				GOPATHS = append(GOPATHS, pth)
			}
		}
	}

	pth, err = homedir.Expand("~/go")
	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(pth); err == nil {
		if _, ok = paths[pth]; !ok {
			GOPATHS = append(GOPATHS, pth)
		}
	}

	pth = build.Default.GOPATH
	if _, ok = paths[pth]; !ok {
		GOPATHS = append(GOPATHS, pth)
	}

	GOPATH = GOPATHS[0]
}

func ResolveGoPath(pth string) (gopath string) {
	for _, gopath := range GOPATHS {
		gpth := path.Join(gopath, pth)
		if _, err := os.Stat(gpth); err == nil {
			return gopath
		}
	}
	return ""
}

func ResolveGoSrcPath(pth string) string {
	pth = path.Join("src", pth)
	for _, gopath := range GOPATHS {
		gpth := path.Join(gopath, pth)
		if _, err := os.Stat(gpth); err == nil {
			return gpth
		}
	}
	return ""
}

func IsExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}

func IsExistingRegularFile(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsRegular()
	}
	return false
}
