package filepresenter

import (
	"fmt"
	"os"
)

type FilePresenter struct {
	path string
}

func NewFilePresenter(path string) *FilePresenter {
	if path == "" {
		path = "file.txt"
	}
	return &FilePresenter{path: path}
}
func (p *FilePresenter) Present(data []string) error {
	file, err := os.OpenFile(p.path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	for _, d := range data {
		_, err = file.WriteString(d + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
