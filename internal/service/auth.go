package service

type AuthService struct {
	Owner    OwnerAuthService
	User     UserAuthService
	Pharmacy PharmacyAuthService
}

type OwnerAuthService interface {
	Login() error
	Logout() error
}
type UserAuthService interface {
	Login() error
	Logout() error
}
type PharmacyAuthService interface {
	Connect() error
}

type ownerAuthService struct{}

func (s *ownerAuthService) Login() error  { return nil }
func (s *ownerAuthService) Logout() error { return nil }

type userAuthService struct{}

func (s *userAuthService) Login() error  { return nil }
func (s *userAuthService) Logout() error { return nil }

type pharmacyAuthService struct{}

func (s *pharmacyAuthService) Connect() error { return nil }

func NewAuthService() *AuthService {
	return &AuthService{
		Owner:    &ownerAuthService{},
		User:     &userAuthService{},
		Pharmacy: &pharmacyAuthService{},
	}
}
