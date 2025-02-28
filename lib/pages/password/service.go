package password

import "fmt"

type Service struct {
	password string
}

func NewService(password string) *Service {
	return &Service{password: password}
}

func (s *Service) VerifyPassword(password string) error {
	if password != s.password {
		return fmt.Errorf("密码错误")
	}
	return nil
}
