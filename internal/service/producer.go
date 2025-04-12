package service

type Producer interface {
	Produce() ([]string, error)
}
