package service

import (
	"context"
	"log/slog"
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

func (s *Service) Run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := sync.WaitGroup{}
	sizeBuffer := 10
	semChan := make(chan struct{}, sizeBuffer)

	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	dataChan := make(chan string, len(data))
	resultChan := make(chan string, len(data))

	for i := range data {
		dataChan <- data[i]
	}
	close(dataChan)

	wg.Add(sizeBuffer)
	for i := 0; i < sizeBuffer; i++ {
		go func(i int) {
			logger.Debug("Goroutine is starting", "goroutine ID", i)
			semChan <- struct{}{}
			defer func() {
				<-semChan
				logger.Debug("Goroutine is ended", "goroutine ID", i)
				wg.Done()
			}()

			for val := range dataChan {
				select {
				case <-ctx.Done():
					logger.Error("Context was canceled", "Error ctx:", ctx.Err(), "goroutine ID", i)
					return
				default:
					resultChan <- s.masking(val)
					//time.Sleep(1500 * time.Millisecond)
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(semChan)
	}()

	var resultSlice []string
	for result := range resultChan {
		resultSlice = append(resultSlice, result)
	}

	err = s.pres.Present(resultSlice)
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
