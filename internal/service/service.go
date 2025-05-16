package service

import (
	"sync"
)

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{prod, pres}
}

func (s *Service) masking(data string) string {
	arr := []byte(data)
	templateUrl := []byte("https://")
	lenTemplateUrl := len(templateUrl)
	for i := 0; i < len(arr)-1; i++ {
		if i > len(arr)-1-lenTemplateUrl {
			break
		}
		a := string(arr[i : i+lenTemplateUrl])
		b := string(templateUrl)
		if a == b {
			for index := i + lenTemplateUrl; index < len(arr); index++ {
				if string(arr[index]) == " " {
					break
				}
				arr[index] = []byte("*")[0]
			}
		}
	}
	return string(arr)
}

func (s *Service) Run() error {
	wg := sync.WaitGroup{}
	sizeBuffer := 10
	semaphore := NewSemaphore(sizeBuffer)
	resultCh := make(chan string, sizeBuffer)

	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	dataChan := dataToChan(data, sizeBuffer)

	wg.Add(sizeBuffer)
	for i := 0; i < sizeBuffer; i++ {
		go func() {
			semaphore.Acquire()
			defer wg.Done()
			defer semaphore.Release()

			for dataFromChan := range dataChan {
				resultCh <- s.masking(dataFromChan)
			}
		}()
	}

	var resultSlice []string
	collectDone := make(chan struct{})

	go func() {
		for result := range resultCh {
			resultSlice = append(resultSlice, result)
		}
		close(collectDone)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	<-collectDone
		
	err = s.pres.Present(resultSlice)
	if err != nil {
		return err
	}

	return nil
}

func dataToChan(data []string, sizeCh int) chan string {
	resultCh := make(chan string, sizeCh)

	go func() {
		defer close(resultCh)
		for _, v := range data {
			resultCh <- v
		}
	}()

	return resultCh
}
