package fe

import (
	"fmt"
	"os"
)

func GetFileExist(filename string) bool {
	if _, err := os.Stat("names.txt"); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}