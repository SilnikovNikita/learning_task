package main

import (
	"os"
	"task_1/internal/filepresenter"
	"task_1/internal/fileproducer"
	"task_1/internal/service"
)

func InputFile() (string, string) {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("args is empty")
	}
	if len(args) < 2 {
		return args[0], ""

	}
	return args[0], args[1]
}

func main() {
	pathProduce, pathPresent := InputFile()
	producer, err := fileproducer.NewFileProducer(pathProduce)
	if err != nil {
		panic(err)
	}
	presenter := filepresenter.NewFilePresenter(pathPresent)

	newService := service.NewService(producer, presenter)
	err = newService.Run()
	if err != nil {
		panic(err)
	}
}
