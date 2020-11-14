package sys

import (
	"io/ioutil"
	"os"
	"path"
)

// 计算目录的大小
func CalDirSize(dirpath string) (dirsize int64, err error) {
	err = os.Chdir(dirpath)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		dirsize += file.Size()
		if file.Mode().IsDir() {
			subDirSize, err2 := CalDirSize(path.Join(dirpath, file.Name()))
			if err2 == nil {
				dirsize += subDirSize
			}
		}
	}
	return
}
