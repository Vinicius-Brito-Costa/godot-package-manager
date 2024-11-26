package util

import (
	"bufio"
	"errors"
	"os"
)

func GetFile(path string) (*os.File, error){
	if len(path) == 0 {
		return nil, errors.New("Cannot load file from a blank path.")
	}
	Info("Getting file in path: " + path)
	file, err := os.Open(path)
	if err != nil {
		Error("Error trying to load file in path: " + path, err)
		return nil, err
	}
	return file, nil
}

func GetFileLines(file *os.File) ([]string, error) {

	if file == nil {
		Info("Cannot get lines from a nil file.")
		return nil, errors.New("Nil File")
	}

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Info("Successfuly read lines from file.")
	return lines, nil
}
