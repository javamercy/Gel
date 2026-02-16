package core

type SymbolicRefService struct {
	refService *RefService
}

func NewSymbolicRefService(refService *RefService) *SymbolicRefService {
	return &SymbolicRefService{
		refService: refService,
	}
}

func (s *SymbolicRefService) Read(name string, short bool) (string, error) {
	ref, err := s.refService.ReadSymbolic(name)
	if err != nil {
		return "", err
	}
	if short {
		return ref[len("refs/heads/"):], nil
	}
	return ref, nil
}

func (s *SymbolicRefService) Write(name, ref string) error {
	return s.refService.WriteSymbolic(name, ref)
}
