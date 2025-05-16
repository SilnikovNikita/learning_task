package service

type Semaphore struct {
	semaChan chan struct{}
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{
		semaChan: make(chan struct{}, size),
	}
}

func (s *Semaphore) Acquire() {
	s.semaChan <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaChan
}
