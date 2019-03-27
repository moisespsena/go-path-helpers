package path_helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	Real  string
	Alias string
	Data  []byte
}

func CopyTree(dest string, sources []interface{}) (err error) {
	if !IsExistingDir(dest) {
		err = os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for i := len(sources) - 1; i >= 0; i-- {
		var src *Path
		switch s := sources[i].(type) {
		case *Path:
			src = s
		case string:
			src = &Path{Real: s}
		default:
			return fmt.Errorf("Invalid source[%v]: %v", i, sources[i])
		}

		if src.Real == "" {
			p := filepath.Join(dest, src.Alias)
			if !IsExistingDir(filepath.Dir(p)) {
				err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
				if err != nil {
					return err
				}
			}
			f, err := os.Create(p)
			if err != nil {
				return err
			}
			_, err = f.Write(src.Data)
			if err != nil {
				return err
			}
		} else {
			err = filepath.Walk(src.Real, func(path string, info os.FileInfo, err error) error {
				if err == nil {
					var relativePath = strings.TrimPrefix(strings.TrimPrefix(path, src.Real), "/")
					if src.Alias != "" {
						relativePath = filepath.Join(src.Alias, relativePath)
					}
					if info.IsDir() {
						err = os.MkdirAll(filepath.Join(dest, relativePath), os.ModePerm)
					} else if info.Mode().IsRegular() {
						source, err := ioutil.ReadFile(path)
						if err != nil {
							return err
						}
						f, err := os.Create(filepath.Join(dest, relativePath))
						if err != nil {
							return err
						}
						_, err = f.Write(source)
						return err
					}
				}
				return err
			})
		}
		if err != nil {
			return err
		}
	}
	return
}
