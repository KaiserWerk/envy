package common

import "os"

func GetWorkingDir() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}

	return d
}
