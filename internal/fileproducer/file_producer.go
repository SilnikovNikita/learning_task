package fileproducer

import (
	"bufio"
	"fmt"
	"os"
)

type FileProducer struct {
	path string
}

func NewFileProducer(path string) (*FileProducer, error) {
	if path == "" {
		return nil, fmt.Errorf("file producer path is empty")
	}
	return &FileProducer{path: path}, nil
}

func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open(fp.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
