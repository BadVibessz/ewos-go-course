package utils

import (
	dom "hw-async/domain"
	"os"
	"strconv"
)

func CandlePeriodToInt(per dom.CandlePeriod) (int, error) {
	return strconv.Atoi(string(per[:len(per)-1]))
}

func WriteToFile(path string, b []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
