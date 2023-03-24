package files

import (
	"os"

	"github.com/src/main/app/log"
)

func Exist(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		log.Warnf("path: [%s] is dir", path)
		return false
	}

	return true
}
