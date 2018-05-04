package path_helpers

import "os"

func HasSources(skip ...int) bool {
	if len(skip) == 0 || skip[0] == 0 {
		skip = []int{2}
	}
	_, err := os.Stat(GetCalledFileNameSkip(skip[0], true))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return true
}
