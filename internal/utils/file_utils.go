package utils

import "os"

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func CreateDirIfNotExists(path string) error {
	if !IsDir(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func CreateFileIfNotExists(path string) error {
	if !IsFileExists(path) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
