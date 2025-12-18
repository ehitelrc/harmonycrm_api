package services

import (
	"harmony_api/repository"
)

type CaseMVService struct {
	repo *repository.MessageRepository
}

func NewCaseMVService() *CaseMVService {
	return &CaseMVService{
		repo: &repository.MessageRepository{},
	}
}

// func (s *CaseMVService) RefreshAsync() {
// 	go func() {
// 		err := s.repo.RefreshCasesMV(true)
// 		if err != nil {
// 			fmt.Println("❌ Error refreshing MV:", err)
// 		} else {
// 			fmt.Println("✅ MV refreshed successfully")
// 		}
// 	}()
// }
