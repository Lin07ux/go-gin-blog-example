package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	
	return len(content), err
}

func GetExt(filename string) string {
	return path.Ext(filename)
}

func IsNotExist(src string) bool {
	_, err := os.Stat(src)
	
	return os.IsNotExist(err)
}

func IsPermissionDenied(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func MakeDirIfNotExist(src string) error {
	if IsNotExist(src) {
		return MakeDir(src)
	}

	return nil
}

func MakeDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
