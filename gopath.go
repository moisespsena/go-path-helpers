package path_helpers

import (
	"os"
	"fmt"
	"path"
	"strings"
	"go/build"
	"path/filepath"
	"github.com/phayes/permbits"
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

func ResolvPerms(pth string) (perms permbits.PermissionBits, err error) {
	var (
		perms2 permbits.PermissionBits
		err2 error
	)

	pth, err = filepath.Abs(pth)
	if err != nil {
		return
	}

	for {
		perms2, err2 = permbits.Stat(pth)
		if err2 == nil {
			return perms2, nil
		} else if os.IsNotExist(err2) {
			pth = filepath.Dir(pth)
		} else {
			return perms, fmt.Errorf("Fail to get stat of %q: %v", pth, err2)
		}
	}
}

func ResolvFilePerms(pth string) (perms permbits.PermissionBits, err error) {
	if IsExistingRegularFile(pth) {
		return permbits.Stat(pth)
	}

	pth, err = filepath.Abs(pth)

	if err != nil {
		return
	}

	p, err2 := ResolvPerms(filepath.Dir(pth))

	if err2 != nil {
		return perms, err2
	}
	// default file don't have execution perms
	p.SetGroupExecute(false)
	p.SetUserExecute(false)
	p.SetOtherExecute(false)
	return p, nil
}