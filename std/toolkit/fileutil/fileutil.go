package fileutil

import (
	"fmt"
	"github.com/PKUJohnson/solar/std"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
	}
}

func ExistsOrCreateDir(path string) error {
	exist, err := PathExists(path)
	if (exist && err == nil) {
		return nil
	}
	std.LogInfoc("save_file", fmt.Sprintf("Try to make dir %s",path))
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		std.LogErrorc("save_file",err,"Create path error")
		return err
	} else {
		return nil
	}
}
