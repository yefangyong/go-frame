package demo

type Service struct {
	repository *Repository
}

func NewService() *Service {
	repository := NewRepository()
	return &Service{
		repository: repository,
	}
}

// 获取用户ID
func (s *Service) getUserId() []int {
	return s.repository.GetUserIds()
}

// 获取用户
func (s *Service) getUser() []UserModel {
	userIDs := s.repository.GetUserIds()
	return s.repository.GetUserByIds(userIDs)
}
