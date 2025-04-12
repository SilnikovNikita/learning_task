package service

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{prod, pres}
}

func (s *Service) masking(arr []byte) {
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
}

func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	for i, _ := range data {
		bytes := []byte(data[i])
		s.masking(bytes)
		data[i] = string(bytes)
	}

	err = s.pres.Present(data)
	if err != nil {
		return err
	}

	return nil
}
