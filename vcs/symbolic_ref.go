package vcs

type SymbolicRefService struct {
	refService *RefService
}

func NewSymbolicRefService(refService *RefService) *SymbolicRefService {
	return &SymbolicRefService{
		refService: refService,
	}
}

func (s *SymbolicRefService) Read(name string) (string, error) {
	return s.refService.ReadSymbolic(name)
}

func (s *SymbolicRefService) Write(name, ref string) error {
	return s.refService.WriteSymbolic(name, ref)
}
