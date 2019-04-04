package path_helpers

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/moisespsena-go/error-wrap"

	"github.com/mitchellh/go-homedir"
	"github.com/phayes/permbits"
)

var (
	GOPATHC string
	GOPATH  string
	GOPATHS []string
)

func init() {
	var (
		err error
		pth string
		ok  bool
	)

	paths := make(map[string]interface{})

	if _, err := os.Stat("vendor"); err == nil {
		if abs, err := filepath.Abs("vendor"); err == nil {
			paths[abs] = nil
			GOPATHS = append(GOPATHS, abs)
		}
	}

	for _, pth = range strings.Split(os.Getenv("GOPATH"), ":") {
		if pth != "" {
			if _, ok = paths[pth]; !ok {
				GOPATHS = append(GOPATHS, pth)
				paths[pth] = nil
			}
		}
	}

	pth, err = homedir.Expand("~/go")
	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(pth); err == nil {
		if _, ok = paths[pth]; !ok {
			paths[pth] = nil
			GOPATHS = append(GOPATHS, pth)
		}
	}

	pth = build.Default.GOPATH
	if _, ok = paths[pth]; !ok {
		GOPATHS = append(GOPATHS, pth)
		paths[pth] = nil
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

func ResolveGoSrcPath(p ...string) string {
	pth := path.Join("src", path.Join(p...))
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

func IsExistingDirE(pth string) (ok bool, err error) {
	if fi, err := os.Stat(pth); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else if !fi.IsDir() {
		err = fmt.Errorf("%q is not directory", pth)
	} else {
		ok = true
	}
	return
}

func MkdirAll(pth string) error {
	perms, err := ResolvePerms(pth)
	if err != nil {
		return err
	}
	return os.MkdirAll(pth, os.FileMode(perms))
}

func MkdirAllIfNotExists(pth string) error {
	if exists, err := IsExistingDirE(pth); err != nil {
		return err
	} else if exists {
		return nil
	}
	perms, err := ResolvePerms(pth)
	if err != nil {
		return err
	}
	return os.MkdirAll(pth, os.FileMode(perms))
}

func IsExistingRegularFile(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsRegular()
	}
	return false
}

func ResolveMode(pth string) (mode os.FileMode, err error) {
	var perms permbits.PermissionBits
	if perms, err = ResolvePerms(pth); err != nil {
		err = errwrap.Wrap(err, "Resolver permissions of %q", pth)
		return
	}
	mode = os.FileMode(perms)
	return
}

func ResolvePerms(pth string) (perms permbits.PermissionBits, err error) {
	var (
		perms2 permbits.PermissionBits
		err2   error
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

func ResolveFileMode(pth string) (mode os.FileMode, err error) {
	var perms permbits.PermissionBits
	if perms, err = ResolveFilePerms(pth); err != nil {
		err = errwrap.Wrap(err, "Resolver permissions of %q", pth)
		return
	}
	mode = os.FileMode(perms)
	return
}

func ResolveFilePerms(pth string) (perms permbits.PermissionBits, err error) {
	if IsExistingRegularFile(pth) {
		return permbits.Stat(pth)
	}

	pth, err = filepath.Abs(pth)

	if err != nil {
		return
	}

	p, err2 := ResolvePerms(filepath.Dir(pth))

	if err2 != nil {
		return perms, err2
	}
	// default file don't have execution perms
	p.SetGroupExecute(false)
	p.SetUserExecute(false)
	p.SetOtherExecute(false)
	return p, nil
}

func TrimGoPathC(pth string, sub ...string) string {
	if GOPATHC != "" {
		gopathc := filepath.Join(append([]string{GOPATHC}, sub...)...)
		return strings.TrimPrefix(strings.Trim(pth, string(filepath.Separator)),
			strings.TrimPrefix(gopathc, string(filepath.Separator))+string(filepath.Separator))
	}
	return pth
}
