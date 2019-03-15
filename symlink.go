package path_helpers

import "os"

func IsSymlink(mode os.FileMode) bool {
	return mode&os.ModeSymlink == os.ModeSymlink
}
